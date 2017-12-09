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

package watch

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
)

// knownProlbmaticFiles is a list of files that are going to show up in any
// project that are artifacts of the way `bazel query` works. Their absence
// from the filesystem is not an indication that something is wrong, but instead
// an indication that the user has not overridden the system defaults.
var knownProlbmaticFiles = []string{
	"tools/defaults/BUILD",
}

func knownProblematicFile(file string) bool {
	for _, knownProlbmaticFile := range knownProlbmaticFiles {
		if file == knownProlbmaticFile {
			return true
		}
	}
	return false
}

// fsnotify also triggers for file stat and read operations. Explicitly filter the modifying events
// to avoid triggering builds on file acccesses (e.g. due to your IDE checking modified status).
const modifyingEvents = fsnotify.Write | fsnotify.Create | fsnotify.Rename | fsnotify.Remove

type FSNotifyWatcher struct {
	buildFileWatcher  *fsnotify.Watcher
	sourceFileWatcher *fsnotify.Watcher

	filesWatched       map[*fsnotify.Watcher]map[string]bool // Inner map is a surrogate for a set
	sourceEventHandler *SourceEventHandler

	buildEvents  chan fsnotify.Event
	sourceEvents chan fsnotify.Event
}

/**
 * NewFSNotifyWatcher creates a new FSNotifyWatcher which tracks changes for
 * build/source files and emits events on a channel when they have been
 * modified.
 */
func NewFSNotifyWatcher() (*FSNotifyWatcher, error) {
	r := &FSNotifyWatcher{}

	r.filesWatched = map[*fsnotify.Watcher]map[string]bool{}

	var err error
	// Even though we are going to recreate this when the query happens, create
	// the pointer we will use to refer to the watchers right now.
	r.buildFileWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	r.buildEvents = make(chan fsnotify.Event, 100)
	go func() {
		for {
			v := <-r.buildFileWatcher.Events
			if v.Op&modifyingEvents != 0 {
				r.buildEvents <- v
			}
		}
	}()

	r.sourceFileWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	r.sourceEvents = make(chan fsnotify.Event, 100)
	go func() {
		for {
			v := <-r.sourceFileWatcher.Events
			if v.Op&modifyingEvents != 0 {
				r.sourceEvents <- v
			}
		}
	}()

	r.sourceEventHandler = NewSourceEventHandler(r.sourceFileWatcher)
	return r, nil
}

func (f *FSNotifyWatcher) WatchSourceFiles(files []string) (count int, err error) {
	return f.watchFiles(files, f.sourceFileWatcher)
}
func (f *FSNotifyWatcher) WatchBuildFiles(files []string) (count int, err error) {
	return f.watchFiles(files, f.buildFileWatcher)
}

func (f *FSNotifyWatcher) watchFiles(files []string, watcher *fsnotify.Watcher) (count int, err error) {
	filesAdded := map[string]bool{}

	for _, line := range files {
		err := watcher.Add(line)
		if err != nil {
			if knownProblematicFile(line) {
				filesAdded[line] = true
				continue
			}
			return 0, fmt.Errorf("Error watching file %v\nError: %v\n", line, err)
		} else {
			filesAdded[line] = true
		}
	}

	for line, _ := range f.filesWatched[watcher] {
		_, ok := filesAdded[line]
		if !ok {
			err := watcher.Remove(line)
			if err != nil {
				return 0, fmt.Errorf("Error watching file %v\nError: %v\n", line, err)
			}
		}
	}

	f.filesWatched[watcher] = filesAdded
	return len(filesAdded), nil
}

func (f *FSNotifyWatcher) Cleanup() {
	f.buildFileWatcher.Close()
	f.sourceFileWatcher.Close()
}

func (f *FSNotifyWatcher) BuildEvents() chan fsnotify.Event {
	return f.buildEvents
}
func (f *FSNotifyWatcher) SourceEvents() chan fsnotify.Event {
	return f.sourceEvents
}
