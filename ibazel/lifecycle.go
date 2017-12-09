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
	"github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf"
)

// Lifecycle is an object that listens to the lifecycle events of iBazel and
// behaves appropriately..
type Lifecycle interface {
	// Initialize is called once it is known that this lifecycle client is going
	// to be used.
	Initialize(info *map[string]string)

	// TargetDecider takes a protobuf rule and performs setup if it matches the
	// listener's expectations.
	TargetDecider(rule *blaze_query.Rule)

	// ChangeDetected is called when a change is detected
	// changeType: "source"|"graph"
	ChangeDetected(targets []string, changeType string, change string)

	// Cleanup is your opportunity to clean up open sockets or connections.
	Cleanup()

	// BeforeCommand is called before a blaze $COMMAND is run.
	// command: "build"|"test"|"run"
	BeforeCommand(targets []string, command string)

	// AfterCommand is called after a blaze $COMMAND is run with the result of
	// that command.
	// command: "build"|"test"|"run"
	AfterCommand(targets []string, command string, success bool)
}
