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
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bazelbuild/bazel-watcher/bazel"
	"github.com/fsnotify/fsnotify"
)

type State string

const (
	DEBOUNCE_QUERY State = "DEBOUNCE_QUERY"
	QUERY          State = "QUERY"
	WAIT           State = "WAIT"
	DEBOUNCE_RUN   State = "DEBOUNCE_RUN"
	RUN            State = "RUN"
	QUIT           State = "QUIT"
)

const debounceDuration = 100 * time.Millisecond
const sourceQuery = "kind('source file', deps(set(%s)))"
const buildQuery = "buildfiles(deps(set(%s)))"

func usage() {
	fmt.Printf(`ibazel

A file watcher for Bazel. Whenever a source file used in a specified
target, run, build, or test the specified targets.

Usage:

ibazel build|test|run targets...

Example:

ibazel test //path/to/my/testing:target
ibazel test //path/to/my/testing/targets/...
ibazel run //path/to/my/runnable:target
ibazel build //path/to/my/buildable:target

`)
}

func main() {

	if len(os.Args) < 3 {
		usage()
		return
	}

	command := os.Args[1]
	targets := os.Args[2:]

	// Even though we are going to recreate this when the query happens, create
	// the pointer we will use to refer to the watchers right now.
	buildFileWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Watcher error: %v", err)
		return
	}
	defer buildFileWatcher.Close()

	sourceFileWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Watcher error: %v", err)
		return
	}
	defer sourceFileWatcher.Close()

	sourceEventHandler := NewSourceEventHandler(sourceFileWatcher)

	var commandToRun func(string)
	switch command {
	case "build":
		commandToRun = build
	case "test":
		commandToRun = test
	case "run":
		commandToRun = run
	default:
		fmt.Printf("Asked me to perform %s. I don't know how to do that.", command)
		return
	}

	state := QUERY
	for {
		fmt.Printf("State: %s\n", state)
		switch state {
		case WAIT:
			select {
			case <-sourceEventHandler.SourceFileEvents:
				fmt.Printf("Detected source change. Rebuilding...\n")
				state = DEBOUNCE_RUN
			case <-buildFileWatcher.Events:
				fmt.Printf("Detected build graph change. Requerying...\n")
				state = DEBOUNCE_QUERY
			}
		case DEBOUNCE_QUERY:
			select {
			case <-buildFileWatcher.Events:
				state = DEBOUNCE_QUERY
			case <-time.After(debounceDuration):
				state = QUERY
			}
		case QUERY:
	                // Query for which files to watch.
                        fmt.Printf("Querying for BUILD files...\n")
                        for _, target := range targets {
				watchFiles(fmt.Sprintf(buildQuery, target), buildFileWatcher)
                        }
                        fmt.Printf("Querying for source files...\n")
                        for _, target := range targets {
				watchFiles(fmt.Sprintf(sourceQuery, target), sourceFileWatcher)
                        }
                        state = RUN
		case DEBOUNCE_RUN:
			select {
			case <-sourceEventHandler.SourceFileEvents:
				state = DEBOUNCE_RUN
			case <-time.After(debounceDuration):
				state = RUN
			}
		case RUN:
			state = WAIT
                        for _, target := range targets {
				fmt.Printf("%sing %s\n", strings.Title(command), target)
				commandToRun(target)
                        }
		}
	}
}

func queryForSourceFiles(query string) []string {
	b := bazel.New()
	b.WriteToStderr(false)
	b.WriteToStdout(false)

	res, err := b.Query(query)
	if err != nil {
		fmt.Printf("Error running Bazel %s\n", err)
	}

	toWatch := make([]string, 0, 10000)
	for _, line := range res {
		if strings.HasPrefix(line, "@") {
			continue
		}
		if strings.HasPrefix(line, "//external") {
			continue
		}

		// For files that are served from the root they will being with "//:". This
		// is a problematic string because, for example, "//:demo.sh" will become
		// "/demo.sh" which is in the root of the filesystem and is unlikely to exist.
		if strings.HasPrefix(line, "//:") {
			line = line[3:]
		}

		toWatch = append(toWatch, strings.Replace(strings.TrimPrefix(line, "//"), ":", "/", 1))
	}

	return toWatch
}

func watchFiles(query string, watcher *fsnotify.Watcher) {
	toWatch := queryForSourceFiles(query)

	// TODO: Figure out how to unwatch files that are no longer included

	for _, line := range toWatch {
		fmt.Printf("Watching: %s\n", line)
		err := watcher.Add(line)
		if err != nil {
			fmt.Printf("Error watching file %v\nError: %v\n", line, err)
			continue
		}
	}
}

func build(target string) {
	b := bazel.New()

	b.Cancel()
	b.WriteToStderr(true)
	b.WriteToStdout(true)
	err := b.Build(target)
	if err != nil {
		fmt.Printf("Build error: %v", err)
		return
	}
}

func test(target string) {
	b := bazel.New()

	b.Cancel()
	b.WriteToStderr(true)
	b.WriteToStdout(true)
	err := b.Test(target)
	if err != nil {
		fmt.Printf("Build error: %v", err)
		return
	}
}

func run(target string) {
	b := bazel.New()

	b.Cancel()
	b.WriteToStderr(true)
	b.WriteToStdout(true)

	// Start run in a goroutine so that it doesn't block watching for files that
	// have changed.
	go b.Run(target)
}
