package main

import (
	"github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf"
)

// Lifecycle is an object that listens to the lifecycle events of iBazel and
// behaves appropriately..
type Lifecycle interface {
	// Initialize is called once it is known that this lifecycle client is going
	// to be used.
	Initialize()

	// TargetDecider takes a protobuf rule and performs setup if it matches the
	// listener's expectations.
	TargetDecider(rule *blaze_query.Rule)

	// ChangeDetected is called when a change is detected
	// changeType: "source"|"graph"
	ChangeDetected(changeType string)

	// Cleanup is your opportunity to clean up open sockets or connections.
	Cleanup()

	// BeforeCommand is called before a blaze $COMMAND is run.
	// command: "build"|"test"|"run"
	BeforeCommand(command string)
	// AfterCommand is called after a blaze $COMMAND is run with the result of
	// that command.
	// command: "build"|"test"|"run"
	AfterCommand(command string, success bool)
}
