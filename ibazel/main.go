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
)

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

// main entrypoint for IBazel.
func main() {
	if len(os.Args) < 3 {
		usage()
		return
	}

	command := strings.ToLower(os.Args[1])

	i, err := New()
	if err != nil {
		fmt.Printf("Error creating iBazel", err)
		os.Exit(1)
	}
	defer i.Cleanup()

	switch command {
	case "build":
		i.Build(os.Args[2:]...)
	case "test":
		i.Test(os.Args[2:]...)
	case "run":
		i.Run(os.Args[2])
	default:
		fmt.Printf("Asked me to perform %s. I don't know how to do that.", command)
		usage()
		return
	}
}
