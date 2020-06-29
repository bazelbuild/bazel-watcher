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

package fsnotify

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"

	"github.com/bazelbuild/bazel-watcher/ibazel/fswatcher/common"
)

// We have to declare our own partial interface in order to mock it out in test
// as the real struct varies from platform to platform
type fsNotifyWatcher interface {
	Add(name string) error
	Remove(name string) error
	Close() error
	Events() chan fsnotify.Event
}
type fsNotifyWatcherWrapper struct {
	watcher *fsnotify.Watcher
}
func (w *fsNotifyWatcherWrapper) Add(name string) error       { return w.watcher.Add(name) }
func (w *fsNotifyWatcherWrapper) Remove(name string) error    { return w.watcher.Remove(name) }
func (w *fsNotifyWatcherWrapper) Close() error                { return w.watcher.Close() }
func (w *fsNotifyWatcherWrapper) Events() chan fsnotify.Event { return w.watcher.Events }

type realFSNotifyWatcher struct {
	watched map[string]struct{}
	wrapper fsNotifyWatcher
}

var _ common.Watcher = &realFSNotifyWatcher{}

// UpdateAll implements ibazel/fswatcher/common.Watcher
func (w *realFSNotifyWatcher) UpdateAll(names []string) error {
	var errs []string
	prev_watched := w.watched
	new_watched := make(map[string]struct{}, len(names))

	for _, name := range names {
		new_watched[name] = struct{}{}

		_, ok := prev_watched[name]
		if ok {
			delete(w.watched, name)
		} else {
			err := w.wrapper.Add(name)
			if err != nil {
				errs = append(errs, fmt.Sprintf("Error watching file %q error: %v", name, err))
			}
		}
	}

	for name, _ := range prev_watched {
		err := w.wrapper.Remove(name)
		if err != nil {
			errs = append(errs, fmt.Sprintf("Error unwatching file %q error: %v\n", name, err))
		}
	}

	w.watched = new_watched

	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "\n"))
	}
	return nil
}

// Close implements ibazel/fswatcher/common.Watcher
func (w *realFSNotifyWatcher) Close() error {
	return w.wrapper.Close()
}

// Events implements ibazel/fswatcher/common.Watcher
func (w *realFSNotifyWatcher) Events() chan common.Event {
	return w.wrapper.Events()
}

func NewWatcher() (common.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	wrapper := &fsNotifyWatcherWrapper{watcher: watcher}
	return &realFSNotifyWatcher{wrapper: wrapper}, err
}
