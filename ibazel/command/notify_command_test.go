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

package command

import (
	"errors"
	"testing"

	"github.com/bazelbuild/bazel-watcher/bazel"
	mock_bazel "github.com/bazelbuild/bazel-watcher/bazel/testing"
	"github.com/bazelbuild/bazel-watcher/process_group"
)

func TestNotifyCommand(t *testing.T) {
	pg := process_group.Command("cat")

	c := &notifyCommand{
		args:      []string{"moo"},
		bazelArgs: []string{},
		pg:        pg,
		target:    "//path/to:target",
	}

	if c.IsSubprocessRunning() {
		t.Errorf("New subprocess shouldn't have been started yet. State: %v", pg.RootProcess().ProcessState)
	}

	var err error
	c.stdin, err = pg.RootProcess().StdinPipe()
	if err != nil {
		t.Error(err)
	}

	// Mock out bazel to return non-error on test
	b := &mock_bazel.MockBazel{}
	b.BuildError(nil)
	bazelNew = func() bazel.Bazel { return b }
	defer func() { bazelNew = oldBazelNew }()

	c.NotifyOfChanges()
	b.BuildError(errors.New("Demo error"))
	c.NotifyOfChanges()
	b.BuildError(nil)
	c.NotifyOfChanges()

	b.AssertActions(t, [][]string{
		{"WriteToStderr"},
		{"WriteToStdout"},
		{"Build", "//path/to:target"},
		{"WriteToStderr"},
		{"WriteToStdout"},
		{"Run", "--script_path=.*", "//path/to:target"},
		{"WriteToStderr"},
		{"WriteToStdout"},
		{"Build", "//path/to:target"},
		{"WriteToStderr"},
		{"WriteToStdout"},
		{"Build", "//path/to:target"},
		{"WriteToStderr"},
		{"WriteToStdout"},
		{"Run", "--script_path=.*", "//path/to:target"},
	})
}

func TestNotifyCommand_ShortCircuit(t *testing.T) {
	pg := process_group.Command("cat")

	c := &notifyCommand{
		args:         []string{"moo"},
		bazelArgs:    []string{},
		pg:           pg,
		target:       "//path/to:target",
		shortCircuit: true,
	}

	if c.IsSubprocessRunning() {
		t.Errorf("New subprocess shouldn't have been started yet. State: %v", pg.RootProcess().ProcessState)
	}

	var err error
	c.stdin, err = pg.RootProcess().StdinPipe()
	if err != nil {
		t.Error(err)
	}

	// Mock out bazel to return non-error on test
	b := &mock_bazel.MockBazel{}
	b.BuildError(nil)
	bazelNew = func() bazel.Bazel { return b }
	defer func() { bazelNew = oldBazelNew }()

	c.NotifyOfChanges()
	b.BuildError(errors.New("Demo error"))
	c.NotifyOfChanges()
	b.BuildError(nil)
	c.NotifyOfChanges()

	b.AssertActions(t, [][]string{
		{"WriteToStderr"},
		{"WriteToStdout"},
		{"Build", "//path/to:target"},
		{"WriteToStderr"},
		{"WriteToStdout"},
		{"Run", "--script_path=.*", "//path/to:target"},
		{"Cancel"},
		{"WriteToStderr"},
		{"WriteToStdout"},
		{"Build", "//path/to:target"},
		{"Cancel"},
		{"WriteToStderr"},
		{"WriteToStdout"},
		{"Build", "//path/to:target"},
		{"WriteToStderr"},
		{"WriteToStdout"},
		{"Run", "--script_path=.*", "//path/to:target"},
	})
}

func TestNotifyCommand_Restart(t *testing.T) {
	var pg process_group.ProcessGroup

	pg = process_group.Command("ls")
	execCommand = func(name string, args ...string) process_group.ProcessGroup {
		return oldExecCommand("ls")
	}
	defer func() { execCommand = oldExecCommand }()

	c := &notifyCommand{
		args:      []string{"moo"},
		bazelArgs: []string{},
		pg:        pg,
		target:    "//path/to:target",
	}

	var err error
	c.stdin, err = pg.RootProcess().StdinPipe()
	if err != nil {
		t.Error(err)
	}

	b := &mock_bazel.MockBazel{}
	b.BuildError(errors.New("Demo error"))
	bazelNew = func() bazel.Bazel { return b }
	defer func() { bazelNew = oldBazelNew }()

	if c.IsSubprocessRunning() {
		t.Errorf("new subprocess shouldn't have been started yet. State: %v", pg.RootProcess().ProcessState)
	}

	c.NotifyOfChanges()
	if c.IsSubprocessRunning() {
		t.Errorf("process should not start with build errors. State: %v", pg.RootProcess().ProcessState)
	}

	// Since the process isn't currently running, this should start it.
	b.BuildError(nil)
	c.NotifyOfChanges()
	if !c.IsSubprocessRunning() {
		t.Errorf("subprocess should have started. State: %v", pg.RootProcess().ProcessState)
	}

	pid1 := c.pg.RootProcess().Process.Pid

	c.Terminate()
	if c.IsSubprocessRunning() {
		t.Errorf("subprocess should have been terminated. State: %v", pg.RootProcess().ProcessState)
	}

	b.BuildError(errors.New("Demo error"))
	c.NotifyOfChanges()
	if c.IsSubprocessRunning() {
		t.Errorf("subprocess should not restart with build errors. State: %v", pg.RootProcess().ProcessState)
	}

	// Since the process isn't currently running, this should re-start it.
	b.BuildError(nil)
	c.NotifyOfChanges()
	if !c.IsSubprocessRunning() {
		t.Errorf("subprocess should have been restarted. State: %v", pg.RootProcess().ProcessState)
	}

	pid2 := c.pg.RootProcess().Process.Pid
	if pid2 == pid1 {
		t.Error("PIDs of restarted process should be different that original process")
	}

	c.NotifyOfChanges()
	if pid2 != c.pg.RootProcess().Process.Pid {
		t.Error("non-dead process was restarted")
	}
}
