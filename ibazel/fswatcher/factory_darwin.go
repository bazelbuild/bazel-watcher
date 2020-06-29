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

// +build darwin

package fswatcher

import (
	"os"
	"sync"

	"github.com/bazelbuild/bazel-watcher/ibazel/fswatcher/fsevents"
	"github.com/bazelbuild/bazel-watcher/ibazel/fswatcher/fsnotify"
	"github.com/bazelbuild/bazel-watcher/ibazel/log"
)

var experimentalWatcherLog sync.Once

func NewWatcher() (Watcher, error) {
	flag, ok := os.LookupEnv("IBAZEL_USE_LEGACY_WATCHER")
	if ok && flag != "0" {
		return fsnotify.NewWatcher()
	}
	experimentalWatcherLog.Do(func() {
		log.Log("You are using an experimental filesystem watcher. If you would like to disable that, please set the environment variable\n\tIBAZEL_USE_LEGACY_WATCHER=1")
	})
	return fsevents.NewWatcher()
}
