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

var b *bazel.Bazel

const moveCheckInterval time.Duration = 20 * time.Millisecond

const moveCheckRetries int = 10

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

	b = bazel.New()

	command := os.Args[1]
	target := os.Args[2]

	query := fmt.Sprintf("kind('source file', deps('%s'))", target)

	toWatch := queryForSourceFiles(query)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Watcher error: %v", err)
		return
	}
	defer watcher.Close()

	for _, line := range toWatch {
		fmt.Printf("Line: %s\n", line)
		err = watcher.Add(line)
		if err != nil {
			fmt.Printf("Error watching: %v", err)
			return
		}
	}

	var commandToRun func(*bazel.Bazel, string)
	switch command {
	case "build":
		fmt.Printf("Building %s\n", target)
		commandToRun = build
	case "test":
		fmt.Printf("Testing %s\n", target)
		commandToRun = test
	case "run":
		fmt.Printf("Running %s\n", target)
		commandToRun = run
	default:
		fmt.Printf("Asked me to perform %s. I don't know how to do that.", command)
		return
	}

	// Listen to the events and trigger action based on the response code.
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				switch event.Op {
				case fsnotify.Remove:
					for i := 0; i < moveCheckRetries; i++ {
						err = watcher.Add(event.Name)
						if err == nil {
							fmt.Printf("File replaced, rebuilding...\n")
							break
						}
						if i == moveCheckRetries - 1 {
							fmt.Printf("File removed, rebuilding...\n")
						}
						time.Sleep(moveCheckInterval)
					}
				default:
					fmt.Printf("File changed, rebuilding...\n")
				}
				commandToRun(b, target)
			case err := <-watcher.Errors:
				fmt.Println("Error:", err)
			}
		}
	}()

	// Kick things off by sending an event to make the first run happen.
	watcher.Events <- fsnotify.Event{}

	// Wait for the file to change for 24 hours. If it doesn't quit.
	time.Sleep(24 * time.Hour)
}

func queryForSourceFiles(query string) []string {
	res, err := b.Query(query)
	if err != nil {
		fmt.Printf("Error running Bazel %s\n", err)
	}

	toWatch := make([]string, 0, 10000)
	for _, line := range res {
		if strings.HasPrefix(line, "@") {
			continue
		}

		toWatch = append(toWatch, strings.Replace(strings.TrimPrefix(line, "//"), ":", "/", 1))
	}
	return toWatch
}

func build(b *bazel.Bazel, target string) {
	b.Cancel()
	b.WriteToStderr(true)
	b.WriteToStdout(true)
	err := b.Build(target)
	if err != nil {
		fmt.Printf("Build error: %v", err)
		return
	}
}

func test(b *bazel.Bazel, target string) {
	b.Cancel()
	b.WriteToStderr(true)
	b.WriteToStdout(true)
	err := b.Test(target)
	if err != nil {
		fmt.Printf("Build error: %v", err)
		return
	}
}

func run(b *bazel.Bazel, target string) {
	b.Cancel()
	b.WriteToStderr(true)
	b.WriteToStdout(true)
	// Start run in a goroutine so that it doesn't block.
	go b.Run(target)
}
