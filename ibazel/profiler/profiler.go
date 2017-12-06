// Copyright 2017 The Bazel Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package profiler

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"
	"errors"
	"log"
	"sync"

	"github.com/satori/go.uuid"
)

const (

	// DefaultPort is the profiler Server's default server port
	DefaultPort uint16 = 30000
)

type Profiler struct {
	server *http.Server
	f *os.File
	targets []string
	iteration string
	iterationStartTime int64
	iterationBuildStart bool
	iterationReloadTriggered bool
	changes []string
	lock sync.Mutex // guards events
}

type profileEvent struct {
	// common
	Type string `json:"type"`
	Iteration string `json:"iteration"`
	Time int64 `json:"time"`
	Targets []string `json:"targets"`
	Elapsed int64 `json:"elapsed,omitempty"`

	// start event
	IBazelVersion string `json:"iBazelVersion,omitempty"`
	BazelVersion string `json:"bazelVersion,omitempty"`
	MaxHeapSize string `json:"maxHeapSize,omitempty"`
	CommittedHeapSize string `json:"committedHeapSize,omitempty"`

	// change event
	Change string `json:"change,omitempty"`

	// build event
	Changes []string `json:"changes,omitempty"`

	// browser event
	RemoteType string `json:"remoteType,omitempty"`
	RemoteTime int64 `json:"remoteTime,omitempty"`
	RemoteElapsed int64 `json:"remoteElapsed,omitempty"`
	RemoteData string `json:"remoteData,omitempty"`
}

type profilerRemoteEvent struct {
	Type string `json:"type"`
	Time int64 `json:"time"`
	Elapsed int64 `json:"elapsed"`
	Data string `json:"data"`
}

func New(outputPath string, targets []string, info *map[string]string) (*Profiler, error) {
	i := &Profiler{}

	f, err := os.OpenFile(outputPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	i.f = f
	i.targets = targets

	i.iterationBuildStart = true
	i.newIteration()
	i.startEvent(info)

	return i, nil
}

func (i *Profiler) Close() {
	if i.f != nil {
		i.f.Close()
	}
	i.closeServer()
}

func (i *Profiler) BuildStartEvent() {
	i.buildEvent("BUILD_START")
}

func (i *Profiler) BuildDoneEvent() {
	i.buildEvent("BUILD_DONE")
}

func (i *Profiler) BuildFailedEvent() {
	i.buildEvent("BUILD_FAILED")
}

func (i *Profiler) TestStartEvent() {
	i.buildEvent("TEST_START")
}

func (i *Profiler) TestDoneEvent() {
	i.buildEvent("TEST_DONE")
}

func (i *Profiler) TestFailedEvent() {
	i.buildEvent("TEST_FAILED")
}

func (i *Profiler) RunStartEvent() {
	i.buildEvent("RUN_START")
}

func (i *Profiler) RunDoneEvent() {
	i.buildEvent("RUN_DONE")
}

func (i *Profiler) SourceChangeEvent(change string) {
	i.changeEvent("SOURCE_CHANGE", change)
}

func (i *Profiler) GraphChangeEvent(change string) {
	i.changeEvent("GRAPH_CHANGE", change)
}

func (i *Profiler) ReloadTriggeredEvent() {
	i.reloadTriggeredEvent()
}

func (i *Profiler) Listen(port uint16) error {
	if i.server != nil {
		return errors.New("Profiler already listening")
	}

	// Create router
	router := http.NewServeMux()

	// Create server
	i.server = &http.Server{
		Handler: router,
		ErrorLog: log.New(os.Stderr, "[profiler]", 0),
	}
	i.server.Addr = makeAddr(port)

	// Handle profiler events
	router.HandleFunc("/profiler-event", profilerEventHandler(i))

	// Create listener
	l, err := net.Listen("tcp", makeAddr(port))
	if err != nil {
		i.closeServer()
		return err
	}

	fmt.Fprintf(os.Stderr, "[profiler] listening on %s\n", i.server.Addr)
	err = i.server.Serve(l)
	i.closeServer()
	return err
}

func (i *Profiler) startEvent(info *map[string]string) {
	event := profileEvent{}
	event.Type = "IBAZEL_START"
	if info != nil {
		event.IBazelVersion = "" // FIXME: get the iBazel version here
		event.BazelVersion = (*info)["release"]
		event.MaxHeapSize = (*info)["max-heap-size"]
		event.CommittedHeapSize = (*info)["committed-heap-size"]
	}
	i.lock.Lock()
	i.processEvent(&event)
	i.lock.Unlock()
}

func (i *Profiler) buildEvent(eventType string) {
	i.lock.Lock()
	i.buildingIteration()
	event := profileEvent{}
	event.Type = eventType
	event.Changes = i.changes
	i.processEvent(&event)
	i.lock.Unlock()
}

func (i *Profiler) changeEvent(eventType string, change string) {
	i.lock.Lock()
	i.newIteration()
	event := profileEvent{}
	event.Type = eventType
	event.Change = change
	i.processEvent(&event)
	i.changes = append(i.changes, change)
	i.lock.Unlock()
}

func (i *Profiler) reloadTriggeredEvent() {
	i.lock.Lock()
	i.iterationReloadTriggered = true
	event := profileEvent{}
	event.Type = "RELOAD_TRIGGERED"
	i.processEvent(&event)
	i.lock.Unlock()
}

func (i *Profiler) remoteEvent(remoteEvent *profilerRemoteEvent) {
	i.lock.Lock()
	if !i.iterationReloadTriggered {
		fmt.Fprintf(os.Stderr, "Ignoring unexpected remote event\n")
		return
	}
	event := profileEvent{}
	event.Type = "REMOTE_EVENT"
	event.RemoteType = remoteEvent.Type
	event.RemoteTime = remoteEvent.Time
	event.RemoteElapsed = remoteEvent.Elapsed
	event.RemoteData = remoteEvent.Data
	i.processEvent(&event)
	i.lock.Unlock()
}

func (i *Profiler) processEvent(event *profileEvent) {
	if i.f != nil && event != nil {
		// prepare the event
		event.Iteration = i.iteration
		event.Time = makeTimestamp()
		event.Targets = i.targets
		event.Elapsed = event.Time - i.iterationStartTime

		// write the event to the output file
		eventJson, _ := json.Marshal(event)
		eventJson = append(eventJson, 10); // \n
		_, err := i.f.Write(eventJson)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to profile file: %v\n", err)
		}
	}
}

func (i *Profiler) newIteration() {
	if i.iterationBuildStart {
		i.iteration = uuid.NewV4().String()
		i.changes = make([]string, 0, 100)
		i.iterationStartTime = makeTimestamp()
		i.iterationBuildStart = false
		i.iterationReloadTriggered = false
	}
}

func (i *Profiler) buildingIteration() {
	i.iterationBuildStart = true
}

func profilerEventHandler(i *Profiler) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			fmt.Fprintf(os.Stderr, "Profiler invalid request method: %s\n", req.Method)
			rw.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		defer req.Body.Close()
		decoder := json.NewDecoder(req.Body)
		var remoteEvent profilerRemoteEvent
		err := decoder.Decode(&remoteEvent)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to decode profile post data: %v\n", err)
			rw.WriteHeader(http.StatusBadRequest)
			return
		}

		fmt.Fprintf(os.Stderr, "Remote event: %s\n", remoteEvent.Type)
		i.remoteEvent(&remoteEvent)
	}
}

func (i *Profiler) closeServer() {
	if i.server != nil {
		err := i.server.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error closing profiler server: %v\n", err)
		}
		i.server = nil
	}
}

// makeAddr converts uint16(x) to ":x"
func makeAddr(port uint16) string {
	return fmt.Sprintf(":%d", port)
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
