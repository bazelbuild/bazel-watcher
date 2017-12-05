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
	"os"
	"time"

	"github.com/satori/go.uuid"
)

type Profiler struct {
	f *os.File
	targets []string;
	iteration string;
	changes []string;
	lastEventTime int64;
}

type profileEvent struct {
	// common
	Type string `json:"type"`
	Iteration string `json:"iteration"`
	Time int64 `json:"time"`
	Targets []string `json:"targets"`
	Duration int64 `json:"duration,omitempty"`

	// start event
	IBazelVersion string `json:"iBazelVersion,omitempty"`
	BazelVersion string `json:"bazelVersion,omitempty"`
	MaxHeapSize string `json:"maxHeapSize,omitempty"`
	CommittedHeapSize string `json:"committedHeapSize,omitempty"`

	// change event
	Change string `json:"change,omitempty"`

	// build event
	Changes []string `json:"changes,omitempty"`
}

func New(outputPath string, targets []string, info *map[string]string) (*Profiler, error) {
	i := &Profiler{}
	i.newIteration()

	f, err := os.OpenFile(outputPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	i.f = f
	i.targets = targets

	i.startEvent(info)

	return i, nil
}

func (i *Profiler) Close() {
	if i.f != nil {
		i.f.Close()
	}
}

func (i *Profiler) BuildStartEvent() {
	i.buildEvent("BUILD_START", false)
}

func (i *Profiler) BuildDoneEvent() {
	i.buildEvent("BUILD_DONE", true)
}

func (i *Profiler) BuildFailedEvent() {
	i.buildEvent("BUILD_FAILED", true)
}

func (i *Profiler) TestStartEvent() {
	i.buildEvent("TEST_START", false)
}

func (i *Profiler) TestDoneEvent() {
	i.buildEvent("TEST_DONE", true)
}

func (i *Profiler) TestFailedEvent() {
	i.buildEvent("TEST_FAILED", true)
}

func (i *Profiler) RunStartEvent() {
	i.buildEvent("RUN_START", false)
}

func (i *Profiler) RunDoneEvent() {
	i.buildEvent("RUN_DONE", true)
}

func (i *Profiler) SourceChangeEvent(change string) {
	i.changeEvent("SOURCE_CHANGE", change)
}

func (i *Profiler) GraphChangeEvent(change string) {
	i.changeEvent("GRAPH_CHANGE", change)
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
	i.processEvent(&event, false)
}

func (i *Profiler) buildEvent(eventType string, buildDone bool) {
	event := profileEvent{}
	event.Type = eventType
	event.Changes = i.changes
	i.processEvent(&event, buildDone)
	if buildDone {
		i.newIteration();
	}
}

func (i *Profiler) changeEvent(eventType string, change string) {
	event := profileEvent{}
	event.Type = eventType
	event.Change = change
	i.processEvent(&event, false)
	i.changes = append(i.changes, change)
}

func (i *Profiler) processEvent(event *profileEvent, calcDuration bool) {
	if i.f != nil && event != nil {
		// prepare the eventiteration
		event.Iteration = i.iteration
		event.Time = makeTimestamp()
		event.Targets = i.targets
		if calcDuration && i.lastEventTime != 0 {
			event.Duration = event.Time - i.lastEventTime
		}
		i.lastEventTime = event.Time

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
	i.iteration = uuid.NewV4().String()
	i.changes = make([]string, 0, 100)
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}