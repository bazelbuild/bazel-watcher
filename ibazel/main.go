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
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bazelbuild/bazel-watcher/ibazel/log"
)

var Version = "Development"

var overrideableStartupFlags []string = []string{
	"--bazelrc",
}

var overrideableBazelFlags []string = []string{
	"--action_env",
	"--announce_rc",
	"--aspects",
	"--build_tag_filters=",
	"--build_tests_only",
	"--compilation_mode",
	"--compile_one_dependency",
	"--config=",
	"--copt=",
	"--curses=no",
	"--cxxopt",
	"-c",
	"--define=",
	"--dynamic_mode=",
	"--features=",
	"--flaky_test_attempts=",
	"--keep_going",
	"-k",
	"--nocache_test_results",
	"--nostamp",
	"--output_groups=",
	"--override_repository=",
	"--platforms",
	"--repo_env",
	"--runs_per_test=",
	"--run_under=",
	"--show_result=",
	"--stamp",
	"--strategy=",
	"--test_arg=",
	"--test_env=",
	"--test_filter=",
	"--test_output=",
	"--test_tag_filters=",
	"--test_timeout=",
	// Custom Starlark build settings
	// https://docs.bazel.build/versions/master/skylark/config.html#using-build-settings-on-the-command-line
	"--//",
	"--no//",
}

var debounceDuration = flag.Duration("debounce", 100*time.Millisecond, "Debounce duration")
var logToFile = flag.String("log_to_file", "-", "Log iBazel stderr to a file instead of os.Stderr")

func usage() {
	fmt.Fprintf(os.Stderr, `iBazel - Version %s

A file watcher for Bazel. Whenever a source file used in a specified
target, run, build, or test the specified targets.

Usage:

ibazel build|test|run [flags] targets...

Example:

ibazel test //path/to/my/testing:target
ibazel test //path/to/my/testing/targets/...
ibazel run //path/to/my/runnable:target -- --arguments --for_your=binary
ibazel build //path/to/my/buildable:target

Supported Bazel startup flags:
  %s

Supported Bazel command flags:
  %s

To add to this list, edit
https://github.com/bazelbuild/bazel-watcher/blob/master/ibazel/main.go

iBazel flags:
`, Version, strings.Join(overrideableStartupFlags, "\n  "), strings.Join(overrideableBazelFlags, "\n  "))
	flag.PrintDefaults()
}

func isOverrideable(arg string, overrideables []string) bool {
	for _, overrideable := range overrideables {
		if strings.HasPrefix(arg, overrideable) {
			return true
		}
	}
	return false
}

func isOverrideableStartupFlag(arg string) bool {
	return isOverrideable(arg, overrideableStartupFlags)
}

func isOverrideableBazelFlag(arg string) bool {
	return isOverrideable(arg, overrideableBazelFlags)
}

func parseArgs(in []string) (targets, startupArgs, bazelArgs, args []string) {
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

			// Check to see if this startup option or command flag is on the bazel whitelist of flags.
			if isOverrideableStartupFlag(arg) {
				startupArgs = append(startupArgs, arg)
			} else if isOverrideableBazelFlag(arg) {
				bazelArgs = append(bazelArgs, arg)
			} else {
				// If none of those things then it's probably a target.
				targets = append(targets, arg)
			}
		}
	}
	return
}

// main entrypoint for IBazel.
func main() {
	flag.Usage = usage
	flag.Parse()

	if *logToFile != "-" {
		var err error
		logFile, err := os.OpenFile(*logToFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		log.SetWriter(logFile)
	}

	if len(flag.Args()) < 2 {
		usage()
		return
	}

	command := strings.ToLower(flag.Args()[0])
	args := flag.Args()[1:]
	os.Setenv("IBAZEL", "true")

	i, err := New()
	if err != nil {
		log.Fatalf("Error creating iBazel: %s", err)
	}
	i.SetDebounceDuration(*debounceDuration)
	defer i.Cleanup()

	// increase the number of files that this process can
	// have open.
	err = setUlimit()
	if err != nil {
		log.Errorf("error setting higher file descriptor limit for this process: %v", err)
	}

	handle(i, command, args)
}

func handle(i *IBazel, command string, args []string) {
	targets, startupArgs, bazelArgs, args := parseArgs(args)
	i.SetStartupArgs(startupArgs)
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
		fmt.Fprintf(os.Stderr, "Asked me to perform %s. I don't know how to do that.", command)
		usage()
		return
	}
}
