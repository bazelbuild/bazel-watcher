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

var overrideableBazelFlags []string = []string{
	"--test_output=",
}

func usage() {
	fmt.Printf(`ibazel

A file watcher for Bazel. Whenever a source file used in a specified
target, run, build, or test the specified targets.

Usage:

ibazel build|test|run targets...

Example:

ibazel test //path/to/my/testing:target
ibazel test //path/to/my/testing/targets/...
ibazel run //path/to/my/runnable:target -- --arguments --for_your=binary
ibazel build //path/to/my/buildable:target

`)
}

func isOverrideableBazelFlag(arg string) bool {
	for _, overrideable := range overrideableBazelFlags {
		if strings.HasPrefix(arg, overrideable) {
			return true
		}
	}
	return false
}

func parseArgs(in []string) (targets, bazelArgs, args []string) {
	afterDoubleDash := false
	for _, arg := range in {
		if afterDoubleDash {
			// Put it in the extra args section if we are after a double dash.
			args = append(args, arg)
		} else {
			// Check to see if this token is a double dash.
			if arg == "--" {
				afterDoubleDash = true
				continue
			}

			// Check to see if this flag is on the bazel whitelist of flags.
			if isOverrideableBazelFlag(arg) {
				bazelArgs = append(bazelArgs, arg)
				continue
			}

			// If none of those things then it's probably a target.
			targets = append(targets, arg)
		}
	}
	return
}

// main entrypoint for IBazel.
func main() {
	if len(os.Args) < 3 {
		usage()
		return
	}

	command := strings.ToLower(os.Args[1])
	args := os.Args[2:]

	i, err := New()
	if err != nil {
		fmt.Printf("Error creating iBazel", err)
		os.Exit(1)
	}
	defer i.Cleanup()

	handle(i, command, args)
}

func handle(i *IBazel, command string, args []string) {
	targets, bazelArgs, args := parseArgs(args)
	i.SetBazelArgs(bazelArgs)

	switch command {
	case "build":
		i.Build(targets...)
	case "test":
		i.Test(targets...)
	case "run":
		// Run only takes one argument
		i.Run(targets[0], args)
	default:
		fmt.Printf("Asked me to perform %s. I don't know how to do that.", command)
		usage()
		return
	}
}
