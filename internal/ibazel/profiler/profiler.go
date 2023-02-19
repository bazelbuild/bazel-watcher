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
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	golog "log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/bazelbuild/bazel-watcher/internal/ibazel/log"
	"github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf/blaze_query"

	_ "embed"
)

var profileDev = flag.String("profile_dev", "", "Turn on profiling and append report to file")

//go:embed profiler.js
var js []byte

const (

	// DefaultPort is the profiler Server's default server port
	DefaultPort uint16 = 30000
)

type Profiler struct {
	server                   *http.Server
	file                     *os.File
	version                  string
	targets                  []string
	iteration                string
	iterationStartTime       int64
	iterationBuildStart      bool
	iterationReloadTriggered bool
	changes                  []string
	lock                     sync.Mutex // guards events
}

type profileEvent struct {
	// common
	Type      string   `json:"type"`
	Iteration string   `json:"iteration"`
	Time      int64    `json:"time"`
	Targets   []string `json:"targets,omitempty"`
	Elapsed   int64    `json:"elapsed,omitempty"`

	// start event
	IBazelVersion     string `json:"iBazelVersion,omitempty"`
	BazelVersion      string `json:"bazelVersion,omitempty"`
	MaxHeapSize       string `json:"maxHeapSize,omitempty"`
	CommittedHeapSize string `json:"committedHeapSize,omitempty"`

	// change event
	Change string `json:"change,omitempty"`

	// build & reload event
	Changes []string `json:"changes,omitempty"`

	// browser event
	RemoteType    string `json:"remoteType,omitempty"`
	RemoteTime    int64  `json:"remoteTime,omitempty"`
	RemoteElapsed int64  `json:"remoteElapsed,omitempty"`
	RemoteData    string `json:"remoteData,omitempty"`
}

type profilerRemoteEvent struct {
	Type                     string `json:"type"`
	Time                     int64  `json:"time"`
	TimeSinceNavigationStart int64  `json:"timeSinceNavigationStart"`
	Data                     string `json:"data"`
}

func New(version string) *Profiler {
	p := &Profiler{}
	p.version = version
	return p
}

func (i *Profiler) Initialize(info *map[string]string) {
	if *profileDev == "" {
		return
	}

	var err error
	i.file, err = os.OpenFile(*profileDev, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Errorf("Failed to open profile output file: %s", *profileDev)
		return
	}

	log.Errorf("Profile output: %s", *profileDev)

	i.iterationBuildStart = true
	i.newIteration()
	i.startEvent(info)
	i.startProfilerServer()
}

func (i *Profiler) TargetDecider(rule *blaze_query.Rule) {}

func (i *Profiler) ChangeDetected(targets []string, changeType string, change string) {
	if i.file == nil {
		return
	}
	i.targets = targets
	switch changeType {
	case "source":
		i.changeEvent("SOURCE_CHANGE", change)
	case "graph":
		i.changeEvent("GRAPH_CHANGE", change)
	}
}

func (i *Profiler) BeforeCommand(targets []string, command string) {
	if i.file == nil {
		return
	}
	i.targets = targets
	switch command {
	case "build":
		i.buildEvent("BUILD_START")
	case "test":
		i.buildEvent("TEST_START")
	case "run":
		i.buildEvent("RUN_START")
	}
}

func (i *Profiler) AfterCommand(targets []string, command string, success bool, output *bytes.Buffer) {
	if i.file == nil {
		return
	}
	i.targets = targets
	if success {
		switch command {
		case "build":
			i.buildEvent("BUILD_DONE")
		case "test":
			i.buildEvent("TEST_DONE")
		case "run":
			i.buildEvent("RUN_DONE")
		}
	} else {
		switch command {
		case "build":
			i.buildEvent("BUILD_FAILED")
		case "test":
			i.buildEvent("TEST_FAILED")
		case "run":
			i.buildEvent("RUN_FAILED")
		}
	}
}

func (i *Profiler) Cleanup() {
	if i.file != nil {
		i.file.Close()
	}
	i.closeServer()
}

func (i *Profiler) ReloadTriggered(targets []string) {
	if i.file == nil {
		return
	}
	i.targets = targets
	i.reloadTriggeredEvent()
}

func (i *Profiler) startProfilerServer() {
	port := DefaultPort
	for ; port < DefaultPort+100; port++ {
		if testPort(port) {
			go func() {
				err := i.listen(port)
				if err != nil {
					log.Errorf("Profiler server failed to start: %v", err)
				}
			}()
			url := fmt.Sprintf("http://localhost:%d/profiler.js", port)
			os.Setenv("IBAZEL_PROFILER_URL", url)
			return
		}
	}
	log.Log("Could not find open port for profiler server")
}

func (i *Profiler) listen(port uint16) error {
	if i.server != nil {
		return errors.New("Profiler already listening")
	}

	// Create router
	router := http.NewServeMux()

	// Create server
	i.server = &http.Server{
		Handler:  router,
		ErrorLog: golog.New(os.Stderr, "[profiler]", 0),
	}
	i.server.Addr = makeAddr(port)

	// Handle profiler.js requests
	router.HandleFunc("/profiler.js", i.jsHandler)

	// Handle profiler events
	router.HandleFunc("/profiler-event", i.profilerEventHandler)

	// Create listener
	l, err := net.Listen("tcp", makeAddr(port))
	if err != nil {
		i.closeServer()
		return err
	}

	log.Logf("[profiler] listening on %s", i.server.Addr)
	err = i.server.Serve(l)
	i.closeServer()
	return err
}

func (i *Profiler) startEvent(info *map[string]string) {
	event := profileEvent{}
	event.Type = "IBAZEL_START"
	if info != nil {
		event.IBazelVersion = i.version
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
	event.Changes = i.changes
	i.processEvent(&event)
	i.lock.Unlock()
}

func (i *Profiler) remoteEvent(remoteEvent *profilerRemoteEvent) {
	i.lock.Lock()
	if !i.iterationReloadTriggered {
		log.Logf("Ignoring unexpected remote event")
		return
	}
	event := profileEvent{}
	event.Type = "REMOTE_EVENT"
	event.RemoteType = remoteEvent.Type
	event.RemoteTime = remoteEvent.Time
	event.RemoteElapsed = remoteEvent.TimeSinceNavigationStart
	event.RemoteData = remoteEvent.Data
	i.processEvent(&event)
	i.lock.Unlock()
}

func (i *Profiler) processEvent(event *profileEvent) {
	if i.file != nil && event != nil {
		// prepare the event
		event.Iteration = i.iteration
		event.Time = makeTimestamp()
		event.Targets = i.targets
		event.Elapsed = event.Time - i.iterationStartTime

		// write the event to the output file
		eventJson, _ := json.Marshal(event)
		eventJson = append(eventJson, 10) // \n
		_, err := i.file.Write(eventJson)
		if err != nil {
			log.Errorf("Error writing to profile file: %v", err)
		}
	}
}

func (i *Profiler) newIteration() {
	if i.iterationBuildStart {
		i.iteration = randomString(16)
		i.changes = make([]string, 0, 100)
		i.iterationStartTime = makeTimestamp()
		i.iterationBuildStart = false
		i.iterationReloadTriggered = false
	}
}

func (i *Profiler) buildingIteration() {
	i.iterationBuildStart = true
}

func (i *Profiler) jsHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		log.Errorf("profiler.js invalid request method: %s", req.Method)
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	rw.Header().Set("Content-Type", "application/javascript")
	_, err := rw.Write(js)
	if err != nil {
		log.Errorf("Error handling profile.js request: %v", err)
	}
}

func (i *Profiler) profilerEventHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		log.Errorf("Profiler invalid request method: %s", req.Method)
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	defer req.Body.Close()
	decoder := json.NewDecoder(req.Body)
	var remoteEvent profilerRemoteEvent
	err := decoder.Decode(&remoteEvent)
	if err != nil {
		log.Errorf("Failed to decode profile post data: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Logf("Remote event: %s", remoteEvent.Type)
	i.remoteEvent(&remoteEvent)
}

func (i *Profiler) closeServer() {
	if i.server != nil {
		err := i.server.Close()
		if err != nil {
			log.Errorf("Error closing profiler server: %v", err)
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

func testPort(port uint16) bool {
	ln, err := net.Listen("tcp", ":"+strconv.FormatInt(int64(port), 10))
	if err != nil {
		log.Errorf("Error opening port %d: %v", port, err)
		return false
	}

	ln.Close()
	return true
}

const letterBytes = "0123456789abcdef"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var randSrc = rand.NewSource(time.Now().UnixNano())

// Fast random string generator
// See: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-golang
func randomString(n int) string {
	b := make([]byte, n)
	// A randSrc.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, randSrc.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randSrc.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
