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
	"reflect"
	"runtime"
	"runtime/debug"
	"syscall"
	"testing"

	"github.com/bazelbuild/bazel-watcher/bazel"
	mock_bazel "github.com/bazelbuild/bazel-watcher/bazel/testing"
	"github.com/bazelbuild/bazel-watcher/ibazel/command"
	"github.com/fsnotify/fsnotify"

	blaze_query "github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf"
	"github.com/golang/protobuf/proto"
)

func assertEqual(t *testing.T, want, got interface{}, msg string) {
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Wanted %s, got %s. %s", want, got, msg)
		debug.PrintStack()
	}
}

type mockCommand struct {
	bazelArgs []string
	target    string
	args      []string

	notifiedOfChanges bool
	started           bool
	terminated        bool
}

func (m *mockCommand) Start() error {
	if m.started {
		panic("Can't run command twice")
	}
	m.started = true
	return nil
}
func (m *mockCommand) NotifyOfChanges() {
	m.notifiedOfChanges = true
}
func (m *mockCommand) Terminate() {
	if !m.started {
		panic("Terminated before starting")
	}
	m.terminated = true
}
func (m *mockCommand) assertTerminated(t *testing.T) {
	if !m.terminated {
		t.Errorf("Mock command wasn't terminated")
		debug.PrintStack()
	}
}

func (m *mockCommand) IsSubprocessRunning() bool {
	return m.started
}

var mockBazel *mock_bazel.MockBazel

func getMockCommand(i *IBazel) *mockCommand {
	c, ok := i.cmd.(*mockCommand)
	if !ok {
		panic(fmt.Sprintf("Unable to cast i.cmd to a mockCommand. Was: %v", i.cmd))
	}
	return c
}

type fakeWatcher struct {
	buildEvents  chan fsnotify.Event
	sourceEvents chan fsnotify.Event
}

func NewMockWatcher() *fakeWatcher {
	return &fakeWatcher{
		buildEvents:  make(chan fsnotify.Event, 1),
		sourceEvents: make(chan fsnotify.Event, 1),
	}
}

func (m *fakeWatcher) Cleanup()                                               {}
func (m *fakeWatcher) WatchSourceFiles(files []string) (count int, err error) { return 0, nil }
func (m *fakeWatcher) WatchBuildFiles(files []string) (count int, err error)  { return 0, nil }
func (m *fakeWatcher) BuildEvents() chan fsnotify.Event {
	return m.buildEvents
}
func (m *fakeWatcher) buildEvent(fsnotify.Event) {
	m.buildEvents <- fsnotify.Event{Op: fsnotify.Write}
}
func (m *fakeWatcher) SourceEvents() chan fsnotify.Event {
	return m.sourceEvents
}
func (m *fakeWatcher) sourceEvent(fsnotify.Event) {
	m.sourceEvents <- fsnotify.Event{Op: fsnotify.Write}
}

func init() {
	// Replace the bazle object creation function with one that makes my mock.
	bazelNew = func() bazel.Bazel {
		mockBazel = &mock_bazel.MockBazel{}
		mockBazel.AddQueryResponse("//path/to:target", &blaze_query.QueryResult{
			Target: []*blaze_query.Target{
				&blaze_query.Target{
					Type: blaze_query.Target_RULE.Enum(),
					Rule: &blaze_query.Rule{
						Attribute: []*blaze_query.Attribute{
							&blaze_query.Attribute{
								Name: proto.String("name"),
							},
						},
					},
				},
			},
		})
		return mockBazel
	}
	commandDefaultCommand = func(bazelArgs []string, target string, args []string) command.Command {
		// Don't do anything
		return &mockCommand{
			bazelArgs: bazelArgs,
			target:    target,
			args:      args,
		}
	}
}

func newIBazel(t *testing.T) *IBazel {
	i, err := New(
		&FakeWorkspaceFinder{},
	)
	if err != nil {
		t.Errorf("Error creating IBazel: %s", err)
	}

	i.lifecycleListeners = []Lifecycle{}
	return i
}

func TestIBazelLoop(t *testing.T) {
	i := newIBazel(t)

	// Replace the file watching channel with one that has a buffer.
	fakeWatcher := NewMockWatcher()
	i.watcher = fakeWatcher
	defer i.Cleanup()

	// The process for testing this is going to be to emit events to the channels
	// that are associated with these objects and walk the state transition
	// graph.

	// First let's consume all the events from all the channels we care about
	called := false
	command := func(targets ...string) error {
		called = true
		return nil
	}

	i.state = QUERY
	step := func() {
		i.iteration("demo", command, []string{}, "")
	}
	assertRun := func() {
		if called == false {
			_, file, line, _ := runtime.Caller(1) // decorate + log + public function.
			t.Errorf("%s:%v Should have run the provided comand", file, line)
		}
		called = false
	}
	assertState := func(state State) {
		if i.state != state {
			_, file, line, _ := runtime.Caller(1) // decorate + log + public function.
			t.Errorf("%s:%v Expected state to be %s but was %s", file, line, state, i.state)
		}
	}

	// Pretend a fairly normal event chain happens.
	// Start, run the program, write a source file, run, write a build file, run.

	assertState(QUERY)
	step()
	assertState(RUN)
	step() // Actually run the command
	assertRun()
	assertState(WAIT)
	// Source file change.
	fakeWatcher.sourceEvent(fsnotify.Event{})
	step()
	assertState(DEBOUNCE_RUN)
	step()
	// Don't send another event in to test the timer
	assertState(RUN)
	step() // Actually run the command
	assertRun()
	assertState(WAIT)
	// Build file change.
	fakeWatcher.buildEvent(fsnotify.Event{})
	step()
	assertState(DEBOUNCE_QUERY)
	return
	// Don't send another event in to test the timer
	step()
	assertState(QUERY)
	step()
	assertState(RUN)
	step() // Actually run the command
	assertRun()
	assertState(WAIT)
}

func TestIBazelBuild(t *testing.T) {
	i := newIBazel(t)
	defer i.Cleanup()

	i.build("//path/to:target")
	expected := [][]string{
		[]string{"Cancel"},
		[]string{"WriteToStderr"},
		[]string{"WriteToStdout"},
		[]string{"Build", "//path/to:target"},
	}

	mockBazel.AssertActions(t, expected)
}

func TestIBazelTest(t *testing.T) {
	t.Skip()
	i := newIBazel(t)
	defer i.Cleanup()

	i.test("//path/to:target")
	expected := [][]string{
		[]string{"Cancel"},
		[]string{"WriteToStderr"},
		[]string{"WriteToStdout"},
		[]string{"Test", "//path/to:target"},
	}

	mockBazel.AssertActions(t, expected)
}

func TestIBazelRun_firstPass(t *testing.T) {
	i := newIBazel(t)
	defer i.Cleanup()

	i.run("//path/to:target")
}

func TestIBazelRun_notifyPreexistiingJobWhenStarting(t *testing.T) {
	commandDefaultCommand = func(bazelArgs []string, target string, args []string) command.Command {
		assertEqual(t, bazelArgs, []string{}, "Bazel args")
		assertEqual(t, target, "", "Target")
		assertEqual(t, args, []string{}, "Args")
		return &mockCommand{}
	}
	defer func() { commandDefaultCommand = command.DefaultCommand }()

	i := newIBazel(t)
	defer i.Cleanup()

	i.args = []string{"--do_it"}

	cmd := &mockCommand{
		notifiedOfChanges: false,
	}
	i.cmd = cmd

	path := "//path/to:target"
	i.run(path)

	if !cmd.notifiedOfChanges {
		t.Errorf("The previously running command was not notified of changes")
	}
}

func TestHandleSignals_SIGINTWithoutRunningCommand(t *testing.T) {
	i := &IBazel{}
	err := i.setup()
	if err != nil {
		t.Errorf("Error creating IBazel: %s", err)
	}
	i.sigs = make(chan os.Signal, 1)
	defer i.Cleanup()

	// But we want to simulate the subprocess not dieing
	attemptedExit := 0
	osExit = func(i int) {
		attemptedExit = i
	}
	assertEqual(t, i.cmd, nil, "There shouldn't be a subprocess running")

	// SIGINT without a running command should attempt to exit
	i.sigs <- syscall.SIGINT
	i.handleSignals()

	// Goroutine tests are kind of racey
	assertEqual(t, attemptedExit, 3, "Should have exited ibazel")
}

func TestHandleSignals_SIGINT(t *testing.T) {
	i := &IBazel{}
	err := i.setup()
	if err != nil {
		t.Errorf("Error creating IBazel: %s", err)
	}
	i.sigs = make(chan os.Signal, 1)
	defer i.Cleanup()

	// But we want to simulate the subprocess not dieing
	attemptedExit := 0
	osExit = func(i int) {
		attemptedExit = i
	}

	var cmd *mockCommand
	// Attempt to kill a task 2 times (but secretly resurrect the job from the
	// dead to test the job not responding)
	for j := 0; j < 2; j++ {
		cmd = &mockCommand{}
		cmd.Start()
		i.cmd = cmd

		// This should kill the subprocess and simulate hitting ctrl-c
		// First save the cmd so we can make assertions on it. It will be removed
		// by the SIGINT
		i.sigs <- syscall.SIGINT
		i.handleSignals()

		cmd.assertTerminated(t)
		assertEqual(t, attemptedExit, 0, "It shouldn't have os.Exit'd")
	}

	// This should kill the job and go over the interrupt limit where exiting happens
	i.sigs <- syscall.SIGINT
	i.handleSignals()

	assertEqual(t, attemptedExit, 3, "Should have exited ibazel")
}

func TestHandleSignals_SIGTERM(t *testing.T) {
	i := &IBazel{}
	err := i.setup()
	if err != nil {
		t.Errorf("Error creating IBazel: %s", err)
	}
	i.sigs = make(chan os.Signal, 1)
	defer i.Cleanup()

	// Now test sending SIGTERM
	attemptedExit := false
	osExit = func(i int) {
		attemptedExit = true
	}
	attemptedExit = false

	cmd := &mockCommand{}
	cmd.Start()
	i.cmd = cmd

	i.sigs <- syscall.SIGTERM
	i.handleSignals()
	cmd.assertTerminated(t)

	assertEqual(t, attemptedExit, true, "Should have exited ibazel")
}
