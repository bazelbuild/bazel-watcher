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
	"os"
	"runtime"
	"syscall"
	"testing"

	mock_bazel "github.com/bazelbuild/bazel-watcher/bazel/testing"
	"github.com/bazelbuild/bazel-watcher/ibazel/process_group"
)

func TestDefaultCommand(t *testing.T) {
	var toKill process_group.ProcessGroup

	if runtime.GOOS == "windows" {
		// TODO(jchw): Remove hardcoded path.
		toKill = process_group.Command("C:\\windows\\system32\\notepad")
	} else {
		toKill = process_group.Command("sleep", "10s")
	}

	execCommand = func(name string, args ...string) process_group.ProcessGroup {
		if runtime.GOOS == "windows" {
			// TODO(jchw): Remove hardcoded path.
			return oldExecCommand("C:\\windows\\system32\\where")
		}
		return oldExecCommand("ls") // Every system has ls.
	}
	defer func() { execCommand = oldExecCommand }()

	c := &defaultCommand{
		args:      []string{"moo"},
		bazelArgs: []string{},
		pg:        toKill,
		target:    "//path/to:target",
		signum:    syscall.SIGKILL,
	}

	if c.IsSubprocessRunning() {
		t.Errorf("New subprocess shouldn't have been started yet. State: %v", toKill.RootProcess().ProcessState)
	}

	toKill.Start()

	if !c.IsSubprocessRunning() {
		t.Errorf("New subprocess was never started. State: %v", toKill.RootProcess().ProcessState)
	}

	// This is synonymous with killing the job so use it to kill the job and test everything.
	c.NotifyOfChanges()
	assertKilled(t, toKill.RootProcess())
}

func TestDefaultCommand_Start(t *testing.T) {
	// Set up mock execCommand and prep it to be returned
	execCommand = func(name string, args ...string) process_group.ProcessGroup {
		if runtime.GOOS == "windows" {
			// TODO(jchw): Remove hardcoded path.
			return oldExecCommand("C:\\windows\\system32\\where")
		}
		return oldExecCommand("ls") // Every system has ls.
	}
	defer func() { execCommand = oldExecCommand }()

	b := &mock_bazel.MockBazel{}

	_, pg := start(b, "//path/to:target", []string{"moo"})
	pg.Start()

	if pg.RootProcess().Stdout != os.Stdout {
		t.Errorf("Didn't set Stdout correctly")
	}
	if pg.RootProcess().Stderr != os.Stderr {
		t.Errorf("Didn't set Stderr correctly")
	}

	b.AssertActions(t, [][]string{
		[]string{"Run", "--script_path=.*", "//path/to:target"},
	})
}
