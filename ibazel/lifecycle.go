package main

import (
	"github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf"
)

// Lifecycle is an object that listens to the lifecycle events of iBazel and
// behaves appropriately..
type Lifecycle interface {
	// TargetDecider takes a protobuf rule and performs setup if it matches the
	// listener's expectations.
	TargetDecider(rule *blaze_query.Rule)

	// Setup is called once it is known that this lifesycle client is going to be
	// used. You can call methods on the ibazel object in this context to set
	// additional environment variables to be passed into the client or to
	// retrigger action.
	Setup()
	// Cleanup is your opportunity to clean up open sockets or connections.
	Cleanup()

	// Before running an "event" where name = (build|test|run).
	BeforeEvent(string)
	AfterEvent(string)
}
