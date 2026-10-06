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

package ibazel

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	"syscall"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/bazelbuild/bazel-watcher/internal/bazel"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/command"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/fswatcher/common"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/log"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/workspace"

	mock_bazel "github.com/bazelbuild/bazel-watcher/internal/bazel/testing"
	analysispb "github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf/analysis"
	blaze_query "github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf/blaze_query"
)

type fakeFSNotifyWatcher struct {
	ErrorChan chan error
	EventChan chan common.Event
}

var _ common.Watcher = &fakeFSNotifyWatcher{}

func (w *fakeFSNotifyWatcher) Close() error                   { return nil }
func (w *fakeFSNotifyWatcher) UpdateAll(names []string) error { return nil }
func (w *fakeFSNotifyWatcher) Events() chan common.Event      { return w.EventChan }

var oldCommandDefaultCommand = command.DefaultCommand

func assertEqual(t *testing.T, want, got interface{}, msg string) {
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Wanted %s, got %s. %s", want, got, msg)
		debug.PrintStack()
	}
}
func assertOsExited(t *testing.T, osExitChan chan int) {
	select {
	case exitCode := <-osExitChan:
		assertEqual(t, 3, exitCode, "Should have exited ibazel")
	case <-time.After(time.Second):
		t.Errorf("It should have os.Exit'd")
		debug.PrintStack()
	}
}
func assertNotOsExited(t *testing.T, osExitChan chan int) {
	select {
	case <-osExitChan:
		t.Errorf("It shouldn't have os.Exit'd")
		debug.PrintStack()
	case <-time.After(time.Second):
		// works as expected
	}
}

type mockCommand struct {
	startupArgs []string
	bazelArgs   []string
	target      string
	args        []string

	notifiedOfChanges bool
	changes           []command.Change
	impact            command.ChangeImpact
	started           bool
	terminated        bool

	signalChan  chan syscall.Signal
	doTermChan  chan struct{}
	didTermChan chan struct{}
}

func (m *mockCommand) Start() (*bytes.Buffer, error) {
	if m.started {
		panic("Can't run command twice")
	}
	m.started = true
	return nil, nil
}
func (m *mockCommand) NotifyOfChanges(changes []command.Change, impact command.ChangeImpact) *bytes.Buffer {
	m.notifiedOfChanges = true
	m.changes = changes
	m.impact = impact
	return nil
}
func (m *mockCommand) Terminate() {
	if !m.started {
		panic("Terminated before starting")
	}
	m.signalChan <- syscall.SIGTERM
	<-m.doTermChan
	m.terminated = true
	m.didTermChan <- struct{}{}
}
func (m *mockCommand) Kill() {
	if !m.started {
		panic("Sending kill signal before terminating")
	}
	m.signalChan <- syscall.SIGKILL
}
func (m *mockCommand) assertTerminated(t *testing.T) {
	select {
	case <-m.didTermChan:
		// works as expected
	case <-time.After(time.Second):
		t.Errorf("A process wasn't terminated within assert timeout")
		debug.PrintStack()
	}
}
func (m *mockCommand) assertSignal(t *testing.T, signum syscall.Signal) {
	if <-m.signalChan != signum {
		t.Errorf("An incorrect signal was used to terminate a process")
		debug.PrintStack()
	}
}
func (m *mockCommand) IsSubprocessRunning() bool {
	return m.started && !m.terminated
}

func getMockCommand(i *IBazel) *mockCommand {
	c, ok := i.cmd.(*mockCommand)
	if !ok {
		panic(fmt.Sprintf("Unable to cast i.cmd to a mockCommand. Was: %v", i.cmd))
	}
	return c
}

func init() {
	commandDefaultCommand = func(startupArgs []string, bazelArgs []string, target string, args []string) command.Command {
		// Don't do anything
		return &mockCommand{
			startupArgs: startupArgs,
			bazelArgs:   bazelArgs,
			target:      target,
			args:        args,
		}
	}
}

func newIBazel(t *testing.T) (*IBazel, *mock_bazel.MockBazel) {
	mockBazel := &mock_bazel.MockBazel{}
	bazelNew = func() bazel.Bazel {
		return mockBazel
	}

	i, err := New("testing")
	if err != nil {
		t.Errorf("Error creating IBazel: %s", err)
	}

	i.workspaceFinder = &workspace.FakeWorkspace{}

	return i, mockBazel
}

func TestIBazelLifecycle(t *testing.T) {
	log.SetTesting(t)

	i, _ := newIBazel(t)
	i.Cleanup()

	// Now inspect private API. If things weren't closed properly this will block
	// and the test will timeout.
	<-i.sourceFileWatcher.Events()
	<-i.buildFileWatcher.Events()
}

func TestIBazelLoop(t *testing.T) {
	log.SetTesting(t)

	i, mockBazel := newIBazel(t)

	ruleType := blaze_query.Target_RULE
	sourceFileType := blaze_query.Target_SOURCE_FILE
	target := "//my:target"
	sourceDir := t.TempDir()

	basePath := filepath.Join(sourceDir, "path", "to")
	if err := os.MkdirAll(basePath, 0o700); err != nil {
		t.Errorf("failed to create based path: %v", err)
		t.FailNow()
	}

	buildFilePath := filepath.Join(basePath, "BUILD")
	sourceFilePath := filepath.Join(basePath, "foo")

	if err := os.WriteFile(buildFilePath, []byte(""), 0o700); err != nil {
		t.Errorf("failed to create BUILD file: %v", err)
		t.FailNow()
	}
	if err := os.WriteFile(sourceFilePath, []byte(""), 0o700); err != nil {
		t.Errorf("failed to create source file: %v", err)
		t.FailNow()
	}

	mockBazel.AddQueryResponse(fmt.Sprintf("buildfiles(set(%q))", target), &blaze_query.QueryResult{
		Target: []*blaze_query.Target{{
			Type: &sourceFileType,
			SourceFile: &blaze_query.SourceFile{
				Name: &buildFilePath,
			},
		}},
	})
	mockBazel.AddCQueryResponse(fmt.Sprintf("deps(set(%s))", target), &analysispb.CqueryResult{
		Results: []*analysispb.ConfiguredTarget{
			{
				Target: &blaze_query.Target{
					Type: &ruleType,
					Rule: &blaze_query.Rule{
						Name:      &target,
						RuleInput: []string{sourceFilePath},
					},
				},
			},
			{
				Target: &blaze_query.Target{
					Type: &sourceFileType,
					SourceFile: &blaze_query.SourceFile{
						Name: &sourceFilePath,
					},
				},
			},
		},
	})

	outputBase := t.TempDir()
	if err := os.Mkdir(filepath.Join(outputBase, "external"), 0o700); err != nil {
		t.Errorf("failed to create external directory: %v", err)
		t.FailNow()
	}
	mockBazel.SetInfo(map[string]string{
		"output_base":  outputBase,
		"install_base": t.TempDir(),
	})

	// Replace the file watching channels with ones that have buffers.
	fakeBuildWatcher := &fakeFSNotifyWatcher{
		EventChan: make(chan common.Event, 1),
	}
	i.buildFileWatcher = fakeBuildWatcher

	fakeSourceWatcher := &fakeFSNotifyWatcher{
		EventChan: make(chan common.Event, 1),
	}
	i.sourceFileWatcher = fakeSourceWatcher

	// Ensure i.filesWatched uses the fake watchers, clearing any potential
	// entries added for the real watchers during i.setup()
	i.filesWatched = make(map[common.Watcher]map[string]struct{})
	i.filesWatched[fakeBuildWatcher] = map[string]struct{}{}
	i.filesWatched[fakeSourceWatcher] = map[string]struct{}{}

	defer i.Cleanup()

	// The process for testing this is going to be to emit events to the channels
	// that are associated with these objects and walk the state transition
	// graph.

	// First let's consume all the events from all the channels we care about
	called := false
	command := func(targets ...string) (*bytes.Buffer, error) {
		called = true
		return nil, nil
	}

	i.state = QUERY
	step := func() {
		i.iteration("demo", command, []string{}, target)
	}
	assertRun := func() {
		t.Helper()

		if called == false {
			_, file, line, _ := runtime.Caller(1) // decorate + log + public function.
			t.Errorf("%s:%v Should have run the provided comand", file, line)
		}
		called = false
	}
	assertState := func(state State) {
		t.Helper()

		if i.state != state {
			_, file, line, _ := runtime.Caller(1) // decorate + log + public function.
			t.Errorf("%s:%v Expected state to be %s but was %s", file, line, state, i.state)
		}
	}

	// Pretend a fairly normal event chain happens.
	// Start, run the program, write a source file, run, write a build file, run.

	assertState(QUERY)
	step()
	i.filesWatched[fakeBuildWatcher][buildFilePath] = struct{}{}
	i.filesWatched[fakeSourceWatcher][sourceFilePath] = struct{}{}
	assertState(RUN)
	step() // Actually run the command
	assertRun()
	assertState(WAIT)
	// Source file change.
	go func() { i.sourceFileWatcher.Events() <- common.Event{Op: common.Write, Name: sourceFilePath} }()
	step()
	assertState(DEBOUNCE_RUN)
	step()
	// Don't send another event in to test the timer
	assertState(RUN)
	step() // Actually run the command
	assertRun()
	assertState(WAIT)
	// Build file change.
	i.buildFileWatcher.Events() <- common.Event{Op: common.Write, Name: buildFilePath}
	step()
	assertState(DEBOUNCE_QUERY)
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
	log.SetTesting(t)

	i, mockBazel := newIBazel(t)
	defer i.Cleanup()

	mockBazel.AddQueryResponse("\"//path/to:target\"", &blaze_query.QueryResult{
		Target: []*blaze_query.Target{
			{
				Type: blaze_query.Target_RULE.Enum(),
				Rule: &blaze_query.Rule{
					Name: proto.String("//path/to:target"),
					Attribute: []*blaze_query.Attribute{
						{Name: proto.String("name")},
					},
				},
			},
		},
	})

	i.build("//path/to:target")
	expected := [][]string{
		{"SetStartupArgs"},
		{"SetArguments"},
		{"Info"},
		{"SetStartupArgs"},
		{"SetArguments"},
		{"Cancel"},
		{"WriteToStderr", "true"},
		{"WriteToStdout", "true"},
		{"Build", "//path/to:target"},
	}

	mockBazel.AssertActions(t, expected)
}

func TestIBazelTest(t *testing.T) {
	log.SetTesting(t)

	i, mockBazel := newIBazel(t)
	defer i.Cleanup()

	mockBazel.AddCQueryResponse("//path/to:target", &analysispb.CqueryResult{
		Results: []*analysispb.ConfiguredTarget{{
			Target: &blaze_query.Target{
				Type: blaze_query.Target_RULE.Enum(),
				Rule: &blaze_query.Rule{
					Name: proto.String("//path/to:target"),
					Attribute: []*blaze_query.Attribute{
						{Name: proto.String("name")},
					},
				},
			},
		}},
	})

	i.test("//path/to:target")
	expected := [][]string{
		{"SetStartupArgs"},
		{"SetArguments"},
		{"Info"},
		{"SetStartupArgs"},
		{"SetArguments"},
		{"SetStartupArgs"},
		{"SetArguments"},
		{"WriteToStderr", "false"},
		{"WriteToStdout", "false"},
		{"CQuery", "//path/to:target"},
		{"SetArguments", "--test_output=streamed"},
		{"Cancel"},
		{"WriteToStderr", "true"},
		{"WriteToStdout", "true"},
		{"Test", "//path/to:target"},
	}

	mockBazel.AssertActions(t, expected)
}

func TestIBazelRun_notifyPreexistiingJobWhenStarting(t *testing.T) {
	log.SetTesting(t)

	commandDefaultCommand = func(startupArgs []string, bazelArgs []string, target string, args []string) command.Command {
		assertEqual(t, startupArgs, []string{}, "Startup args")
		assertEqual(t, bazelArgs, []string{}, "Bazel args")
		assertEqual(t, target, "", "Target")
		assertEqual(t, args, []string{}, "Args")
		return &mockCommand{}
	}
	defer func() { commandDefaultCommand = oldCommandDefaultCommand }()

	i, _ := newIBazel(t)
	defer i.Cleanup()

	i.args = []string{"--do_it"}

	cmd := &mockCommand{
		notifiedOfChanges: false,
	}
	i.cmd = cmd
	i.pendingChanges = []command.Change{{Path: "/workspace/path/to/file", Kind: "source"}}

	path := "//path/to:target"
	i.run(path)

	if !cmd.notifiedOfChanges {
		t.Errorf("The previously running command was not notified of changes")
	}
	assertEqual(t, i.pendingChanges, cmd.changes, "Changes passed to running command")
}

func TestIBazelRunStartsBeforeWatchQuery(t *testing.T) {
	log.SetTesting(t)

	oldCommandNotifyCommand := commandNotifyCommand
	defer func() { commandNotifyCommand = oldCommandNotifyCommand }()

	i, mockBazel := newIBazel(t)
	defer i.Cleanup()

	target := "//path/to:target"
	attributeType := blaze_query.Attribute_STRING_LIST
	mockBazel.AddCQueryResponse(target, &analysispb.CqueryResult{
		Results: []*analysispb.ConfiguredTarget{{
			Target: &blaze_query.Target{
				Type: blaze_query.Target_RULE.Enum(),
				Rule: &blaze_query.Rule{
					Name: proto.String(target),
					Attribute: []*blaze_query.Attribute{{
						Name:            proto.String("tags"),
						Type:            &attributeType,
						StringListValue: []string{"ibazel_notify_changes"},
					}},
				},
			},
		}},
	})
	mockBazel.AddCQueryResponse(fmt.Sprintf("deps(set(%s))", target), &analysispb.CqueryResult{})
	mockBazel.AddQueryResponse("buildfiles(set())", &blaze_query.QueryResult{})
	mockBazel.AddCQueryResponse(fmt.Sprintf("kind('source file', deps(set(%s)))", target), &analysispb.CqueryResult{})

	outputBase := t.TempDir()
	if err := os.Mkdir(filepath.Join(outputBase, "external"), 0o700); err != nil {
		t.Fatalf("failed to create external directory: %v", err)
	}
	mockBazel.SetInfo(map[string]string{
		"output_base":  outputBase,
		"install_base": t.TempDir(),
	})

	cmd := &mockCommand{}
	commandNotifyCommand = func(_ []string, _ []string, _ string, _ []string, _ bool, _ []string) command.Command {
		return cmd
	}

	i.prepareRun(target)
	i.iteration("run", i.run, []string{target}, target)
	if !cmd.started {
		t.Fatal("run target was not started")
	}
	assertEqual(t, QUERY, i.state, "Watch query should run after the target starts")

	i.iteration("run", i.run, []string{target}, target)
	assertEqual(t, RUN, i.state, "Catch-up build should follow the watch query")

	i.iteration("run", i.run, []string{target}, target)
	if !cmd.notifiedOfChanges {
		t.Fatal("run target did not receive the catch-up build")
	}
	assertEqual(t, []command.Change{{Kind: "graph"}}, cmd.changes, "Catch-up build changes")
	assertEqual(t, WAIT, i.state, "Run should wait after the catch-up build")
}

func TestPrepareRunNegotiatesNotificationsAndInitialState(t *testing.T) {
	oldCommandNotifyCommand := commandNotifyCommand
	defer func() { commandNotifyCommand = oldCommandNotifyCommand }()

	for _, test := range []struct {
		name          string
		tags          []string
		structured    bool
		initialState  State
		outputGroups  []string
		directTargets []string
	}{
		{name: "default", initialState: QUERY},
		{name: "legacy", tags: []string{"ibazel_notify_changes"}, initialState: RUN},
		{name: "structured", tags: []string{"ibazel_notify_changes", "ibazel_notify_changes_v1"}, structured: true, initialState: RUN, outputGroups: []string{"generated", "manifest"}, directTargets: []string{"//path/to:rpc", "//path/to:electron"}},
	} {
		t.Run(test.name, func(t *testing.T) {
			i, mockBazel := newIBazel(t)
			defer i.Cleanup()
			i.SetNotifyOutputGroups("generated, manifest")

			target := "//path/to:target"
			attributeType := blaze_query.Attribute_STRING_LIST
			mockBazel.AddCQueryResponse(target, &analysispb.CqueryResult{
				Results: []*analysispb.ConfiguredTarget{{
					Target: &blaze_query.Target{
						Type: blaze_query.Target_RULE.Enum(),
						Rule: &blaze_query.Rule{
							Name:      proto.String(target),
							RuleInput: test.directTargets,
							Attribute: []*blaze_query.Attribute{
								{
									Name:            proto.String("tags"),
									Type:            &attributeType,
									StringListValue: test.tags,
								},
							},
						},
					},
				}},
			})

			var structured bool
			var outputGroups []string
			commandNotifyCommand = func(_ []string, _ []string, _ string, _ []string, enabled bool, groups []string) command.Command {
				structured = enabled
				outputGroups = groups
				return &mockCommand{}
			}

			i.prepareRun(target)
			assertEqual(t, test.structured, structured, "Structured notification mode")
			assertEqual(t, test.initialState, i.state, "Initial run state")
			assertEqual(t, test.outputGroups, outputGroups, "Notification output groups")
			assertEqual(t, test.directTargets, i.directTargets, "Direct targets")
		})
	}
}

func TestChangeDetectedRecordsUniqueChanges(t *testing.T) {
	i := &IBazel{}
	i.changeDetected(nil, "source", "/workspace/path/to/file")
	i.changeDetected(nil, "source", "/workspace/path/to/file")
	i.changeDetected(nil, "graph", "/workspace/path/to/BUILD")

	want := []command.Change{
		{Path: "/workspace/path/to/file", Kind: "source"},
		{Path: "/workspace/path/to/BUILD", Kind: "graph"},
	}
	assertEqual(t, want, i.pendingChanges, "Unique changes for the current iteration")
}

func TestSourceOwnersDriveChangeImpact(t *testing.T) {
	i, mockBazel := newIBazel(t)
	defer i.Cleanup()

	outputBase := t.TempDir()
	if err := os.Mkdir(filepath.Join(outputBase, "external"), 0o700); err != nil {
		t.Fatal(err)
	}
	mockBazel.SetInfo(map[string]string{
		"output_base":  outputBase,
		"install_base": t.TempDir(),
	})

	rpcSource := filepath.Join(t.TempDir(), "rpc.rs")
	electronSource := filepath.Join(t.TempDir(), "electron.ts")
	sharedSource := filepath.Join(t.TempDir(), "shared.txt")
	rpcBuild := filepath.Join(t.TempDir(), "BUILD.bazel")
	for _, path := range []string{rpcSource, electronSource, sharedSource, rpcBuild} {
		if err := os.WriteFile(path, nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}

	ruleType := blaze_query.Target_RULE
	sourceType := blaze_query.Target_SOURCE_FILE
	configuredRule := func(name string, inputs ...string) *analysispb.ConfiguredTarget {
		return &analysispb.ConfiguredTarget{Target: &blaze_query.Target{
			Type: &ruleType,
			Rule: &blaze_query.Rule{Name: proto.String(name), RuleInput: inputs},
		}}
	}
	configuredSource := func(name string) *analysispb.ConfiguredTarget {
		return &analysispb.ConfiguredTarget{Target: &blaze_query.Target{
			Type:       &sourceType,
			SourceFile: &blaze_query.SourceFile{Name: proto.String(name)},
		}}
	}
	rpcRule := configuredRule("//app:rpc", rpcSource, sharedSource)
	rpcRule.Target.Rule.Location = proto.String(rpcBuild + ":1:1")
	graph := &analysispb.CqueryResult{Results: []*analysispb.ConfiguredTarget{
		rpcRule,
		configuredRule("//app:electron", electronSource, sharedSource),
		configuredSource(rpcSource),
		configuredSource(electronSource),
		configuredSource(sharedSource),
	}}

	i.directTargets = []string{"//app:rpc", "//app:electron"}
	i.changeOwners, i.ownershipReady = i.sourceOwnersFromGraph(graph, map[string]string{
		rpcSource:      rpcSource,
		electronSource: electronSource,
		sharedSource:   sharedSource,
	})
	resolved := func(path string) string {
		result, err := filepath.EvalSymlinks(path)
		if err != nil {
			t.Fatal(err)
		}
		return result
	}
	i.pendingChanges = []command.Change{{Path: resolved(rpcSource), Kind: "source"}}
	assertEqual(t, command.ChangeImpact{Targets: []string{"//app:rpc"}, Complete: true}, i.changeImpact(), "RPC change impact")
	i.pendingChanges = []command.Change{{Path: resolved(rpcBuild), Kind: "graph"}}
	assertEqual(t, command.ChangeImpact{Targets: []string{"//app:rpc"}, Complete: true}, i.changeImpact(), "RPC graph impact")

	i.pendingChanges = []command.Change{
		{Path: resolved(sharedSource), Kind: "source"},
		{Path: filepath.Join(t.TempDir(), "BUILD"), Kind: "graph"},
	}
	assertEqual(t, command.ChangeImpact{Targets: []string{"//app:electron", "//app:rpc"}, Complete: false}, i.changeImpact(), "Incomplete shared change impact")
}

func TestHandleSignals_SIGINTWithoutRunningCommand(t *testing.T) {
	log.SetTesting(t)
	log.FakeExit()

	i := &IBazel{}
	err := i.setup()
	if err != nil {
		t.Errorf("Error creating IBazel: %s", err)
	}
	i.sigs = make(chan os.Signal, 1)
	defer i.Cleanup()

	osExitChan := make(chan int, 1)
	osExit = func(i int) {
		osExitChan <- i
	}
	assertEqual(t, i.cmd, nil, "There shouldn't be a subprocess running")

	// SIGINT without a running command should attempt to exit
	i.sigs <- syscall.SIGINT
	i.handleSignals()

	// Goroutine tests are kind of racey
	assertOsExited(t, osExitChan)
}

func TestHandleSignals_SIGINTNormalTermination(t *testing.T) {
	log.SetTesting(t)

	i := &IBazel{}
	err := i.setup()
	if err != nil {
		t.Errorf("Error creating IBazel: %s", err)
	}
	i.sigs = make(chan os.Signal, 1)
	defer i.Cleanup()

	osExitChan := make(chan int, 1)
	osExit = func(i int) {
		osExitChan <- i
	}

	cmd := &mockCommand{
		signalChan:  make(chan syscall.Signal, 10),
		doTermChan:  make(chan struct{}, 1),
		didTermChan: make(chan struct{}, 1),
	}
	i.cmd = cmd
	cmd.Start()

	// First ctrl-c sends custom signal (SIGTERM)
	i.sigs <- syscall.SIGINT
	i.handleSignals()
	cmd.assertSignal(t, syscall.SIGTERM)
	cmd.doTermChan <- struct{}{}
	cmd.assertTerminated(t)
	assertNotOsExited(t, osExitChan)

	// Second ctrl-c terminates ibazel
	i.sigs <- syscall.SIGINT
	i.handleSignals()
	assertOsExited(t, osExitChan)
}

func TestHandleSignals_SIGINTForcefulTermination(t *testing.T) {
	log.SetTesting(t)

	i := &IBazel{}
	err := i.setup()
	if err != nil {
		t.Errorf("Error creating IBazel: %s", err)
	}
	i.sigs = make(chan os.Signal, 1)
	defer i.Cleanup()

	osExitChan := make(chan int, 1)
	osExit = func(i int) {
		osExitChan <- i
	}

	cmd := &mockCommand{
		signalChan:  make(chan syscall.Signal, 10),
		doTermChan:  make(chan struct{}, 1),
		didTermChan: make(chan struct{}, 1),
	}
	i.cmd = cmd
	cmd.Start()

	// First ctrl-c sends custom signal (SIGTERM)
	i.sigs <- syscall.SIGINT
	i.handleSignals()
	cmd.assertSignal(t, syscall.SIGTERM)
	assertNotOsExited(t, osExitChan)

	// Second ctrl-c sends SIGKILL
	i.sigs <- syscall.SIGINT
	i.handleSignals()
	cmd.assertSignal(t, syscall.SIGKILL)
	cmd.doTermChan <- struct{}{}
	cmd.assertTerminated(t)
	assertNotOsExited(t, osExitChan)

	// Yet another ctrl-c terminates ibazel
	i.sigs <- syscall.SIGINT
	i.handleSignals()
	assertOsExited(t, osExitChan)
}

func TestHandleSignals_SIGINTHitLimitTermination(t *testing.T) {
	log.SetTesting(t)

	i := &IBazel{}
	err := i.setup()
	if err != nil {
		t.Errorf("Error creating IBazel: %s", err)
	}
	i.sigs = make(chan os.Signal, 1)
	defer i.Cleanup()

	osExitChan := make(chan int, 1)
	osExit = func(i int) {
		osExitChan <- i
	}

	cmd := &mockCommand{
		signalChan:  make(chan syscall.Signal, 10),
		doTermChan:  make(chan struct{}, 1),
		didTermChan: make(chan struct{}, 1),
	}
	i.cmd = cmd
	cmd.Start()

	// First ctrl-c sends custom signal (SIGTERM)
	i.sigs <- syscall.SIGINT
	i.handleSignals()
	cmd.assertSignal(t, syscall.SIGTERM)
	assertNotOsExited(t, osExitChan)

	// Second ctrl-c sends SIGKILL
	i.sigs <- syscall.SIGINT
	i.handleSignals()
	cmd.assertSignal(t, syscall.SIGKILL)
	assertNotOsExited(t, osExitChan)

	// Third ctrl-c terminates ibazel even if the subprocess is not closed
	i.sigs <- syscall.SIGINT
	i.handleSignals()
	assertOsExited(t, osExitChan)
}

func TestHandleSignals_SIGTERM(t *testing.T) {
	log.SetTesting(t)

	i := &IBazel{}
	err := i.setup()
	if err != nil {
		t.Errorf("Error creating IBazel: %s", err)
	}
	i.sigs = make(chan os.Signal, 1)
	defer i.Cleanup()

	osExitChan := make(chan int, 1)
	osExit = func(i int) {
		osExitChan <- i
	}

	cmd := &mockCommand{
		signalChan:  make(chan syscall.Signal, 10),
		doTermChan:  make(chan struct{}, 1),
		didTermChan: make(chan struct{}, 1),
	}
	i.cmd = cmd
	cmd.Start()

	i.sigs <- syscall.SIGTERM
	i.handleSignals()
	cmd.assertSignal(t, syscall.SIGTERM)
	cmd.doTermChan <- struct{}{}
	cmd.assertTerminated(t)
	assertOsExited(t, osExitChan)
}

func TestParseTarget(t *testing.T) {
	log.SetTesting(t)

	tests := []struct {
		in     string
		repo   string
		target string
	}{
		{"@//my:target", "", "my:target"},
		{"@repo//my:target", "repo", "my:target"},
		{"@bazel_tools//:strange/target", "bazel_tools", ":strange/target"},
	}
	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			gotRepo, gotTarget := parseTarget(test.in)
			if gotRepo != test.repo {
				t.Errorf("parseTarget(%q).repo = %q, want %q", test.in, gotRepo, test.repo)
			}
			if gotTarget != test.target {
				t.Errorf("parseTarget(%q).target = %q, want %q", test.in, gotTarget, test.target)
			}
		})
	}
}

func TestIBazelConvertLabelsToPaths(t *testing.T) {
	log.SetTesting(t)

	i, mockBazel := newIBazel(t)
	defer i.Cleanup()

	outputBase := t.TempDir()
	if err := os.Mkdir(filepath.Join(outputBase, "external"), 0o700); err != nil {
		t.Errorf("failed to create external directory: %v", err)
		t.FailNow()
	}
	mockBazel.SetInfo(map[string]string{
		"output_base":  outputBase,
		"install_base": t.TempDir(),
	})

	repo := t.TempDir()
	if err := os.Symlink(repo, filepath.Join(outputBase, "external", "repo")); err != nil {
		t.Errorf("failed to create symlink to repo: %v", err)
		t.FailNow()
	}

	labels := make([]string, 0, 1)
	labels = append(labels, "//first:file")
	labels = append(labels, "@repo//first:file")
	labels = append(labels, "//second:file")
	labels = append(labels, "@repo//second:file")
	toWatch, err := i.labelsToWatch(labels)
	if err != nil {
		t.Errorf("Error converting labels to paths: %s", err)
	}

	expectedToWatch := make([]string, 0, 1)
	expectedToWatch = append(expectedToWatch, filepath.Join("first", "file"))
	if runtime.GOOS != "windows" {
		expectedToWatch = append(expectedToWatch, filepath.Join(repo, "first", "file"))
	}
	expectedToWatch = append(expectedToWatch, filepath.Join("second", "file"))
	if runtime.GOOS != "windows" {
		expectedToWatch = append(expectedToWatch, filepath.Join(repo, "second", "file"))
	}
	assertEqual(t, expectedToWatch, toWatch, "Resulting paths")
}
