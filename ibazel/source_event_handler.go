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

package main

import (
	"github.com/fsnotify/fsnotify"
)

type SourceEventHandler struct {
	SourceFileEvents  chan fsnotify.Event
	SourceFileWatcher *fsnotify.Watcher
}

func (s *SourceEventHandler) Listen() {
	for {
		select {
		case event := <-s.SourceFileWatcher.Events:
			s.SourceFileEvents <- event

			switch event.Op {
			case fsnotify.Remove, fsnotify.Rename:
				s.SourceFileWatcher.Add(event.Name)
			}
		}
	}
}

func NewSourceEventHandler(sourceFileWatcher *fsnotify.Watcher) *SourceEventHandler {
	handler := &SourceEventHandler{
		make(chan fsnotify.Event),
		sourceFileWatcher,
	}
	go handler.Listen()
	return handler
}
