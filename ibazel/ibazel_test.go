package main

import (
	"os/exec"
	"reflect"
	"runtime"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/bazel"
	"github.com/fsnotify/fsnotify"
)

type MockBazel struct {
	actions       [][]string
	queryResponse []string
}

func (b *MockBazel) WriteToStderr(v bool) {
	b.actions = append(b.actions, []string{"WriteToStderr"})
}
func (b *MockBazel) WriteToStdout(v bool) {
	b.actions = append(b.actions, []string{"WriteToStdout"})
}
func (b *MockBazel) Info() (map[string]string, error) {
	b.actions = append(b.actions, []string{"Info"})
	return map[string]string{}, nil
}
func (b *MockBazel) Query(args ...string) ([]string, error) {
	b.actions = append(b.actions, append([]string{"Query"}, args...))
	return b.queryResponse, nil
}
func (b *MockBazel) Build(args ...string) error {
	b.actions = append(b.actions, append([]string{"Build"}, args...))
	return nil
}
func (b *MockBazel) Test(args ...string) error {
	b.actions = append(b.actions, append([]string{"Test"}, args...))
	return nil
}
func (b *MockBazel) Run(args ...string) (*exec.Cmd, error) {
	b.actions = append(b.actions, append([]string{"Run"}, args...))
	return nil, nil
}
func (b *MockBazel) Cancel() {
	b.actions = append(b.actions, []string{"Cancel"})
}

var mockBazel *MockBazel

func init() {
	// Replace the bazle object creation function with one that makes my mock.
	bazelNew = func() bazel.Bazel {
		mockBazel = &MockBazel{}
		return mockBazel
	}
}

func TestIBazelLifecycle(t *testing.T) {
	i, err := New()
	if err != nil {
		t.Errorf("Error creating IBazel: %s", err)
	}

	i.Cleanup()

	// Now inspect private API. If things weren't closed properly this will block
	// and the test will timeout.
	<-i.sourceFileWatcher.Events
	<-i.buildFileWatcher.Events
}

func TestIBazelLoop(t *testing.T) {
	i, err := New()
	if err != nil {
		t.Errorf("Error creating IBazel: %s", err)
	}

	// Replace the file watching channel with one that has a buffer.
	i.buildFileWatcher.Events = make(chan fsnotify.Event, 1)
	i.sourceEventHandler.SourceFileEvents = make(chan fsnotify.Event, 1)

	defer i.Cleanup()

	// The process for testing this is going to be to emit events to the channels
	// that are associated with these objects and walk the state transition
	// graph.

	// First let's consume all the events from all the channels we care about
	called := false
	command := func(targets ...string) {
		called = true
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
	i.sourceEventHandler.SourceFileEvents <- fsnotify.Event{}
	step()
	assertState(DEBOUNCE_RUN)
	step()
	// Don't send another event in to test the timer
	assertState(RUN)
	step() // Actually run the command
	assertRun()
	assertState(WAIT)
	// Build file change.
	i.buildFileWatcher.Events <- fsnotify.Event{}
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
	i, err := New()
	if err != nil {
		t.Errorf("Error creating IBazel: %s", err)
	}

	defer i.Cleanup()

	i.build("//path/to:target")
	expected := [][]string{
		[]string{"Cancel"},
		[]string{"WriteToStderr"},
		[]string{"WriteToStdout"},
		[]string{"Build", "//path/to:target"},
	}

	if !reflect.DeepEqual(mockBazel.actions, expected) {
		t.Errorf("Build didn't meet expecations.\nWant: %s\nGot:  %s", expected, mockBazel.actions)
	}
}

func TestIBazelTest(t *testing.T) {
	i, err := New()
	if err != nil {
		t.Errorf("Error creating IBazel: %s", err)
	}

	defer i.Cleanup()

	i.test("//path/to:target")
	expected := [][]string{
		[]string{"Cancel"},
		[]string{"WriteToStderr"},
		[]string{"WriteToStdout"},
		[]string{"Test", "//path/to:target"},
	}

	if !reflect.DeepEqual(mockBazel.actions, expected) {
		t.Errorf("Test didn't meet expecations.\nWant: %s\nGot:  %s", expected, mockBazel.actions)
	}
}

func TestIBazelRun(t *testing.T) {
	i, err := New()
	if err != nil {
		t.Errorf("Error creating IBazel: %s", err)
	}

	defer i.Cleanup()

	i.run("//path/to:target")

	expected := [][]string{
		[]string{"Cancel"},
		[]string{"WriteToStderr"},
		[]string{"WriteToStdout"},
		[]string{"Run", "//path/to:target"},
	}

	// Sleep for 50ms to let the goroutine launched by run get scheduled/execute.
	time.Sleep(50 * time.Millisecond)

	if !reflect.DeepEqual(mockBazel.actions, expected) {
		t.Errorf("Test didn't meet expecations.\nWant: %s\nGot:  %s", expected, mockBazel.actions)
	}
}
