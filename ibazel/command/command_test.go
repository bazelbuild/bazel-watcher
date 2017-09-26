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
	"os/exec"
	"testing"

	mock_bazel "github.com/bazelbuild/bazel-watcher/bazel/testing"
)

var oldExecCommand = exec.Command

func assertKilled(t *testing.T, cmd *exec.Cmd) {
	if err := cmd.Wait(); err != nil {
		if cmd.ProcessState.Success() {
			t.Errorf("Subprocess terminated from \"natural\" causes, which means the job ran till its timeout then existed. The Run method should have killed it before then.")
		}
		if cmd.ProcessState == nil {
			t.Errorf("Killable subprocess was never started. State: %v, Err: %v", cmd.ProcessState, err)
		}
	}
}

func TestSubprocessRunning(t *testing.T) {
	if subprocessRunning(nil) {
		t.Errorf("Nil subprocesses don't run")
	}

	cmd := exec.Command("sleep", ".1s")

	if subprocessRunning(cmd) {
		t.Errorf("New subprocess shouldn't have been started yet. State: %v", cmd.ProcessState)
	}

	cmd.Start()

	if !subprocessRunning(cmd) {
		t.Errorf("New subprocess was never started. State: %v", cmd.ProcessState)
	}

	err := cmd.Wait()
	if err != nil || subprocessRunning(cmd) {
		t.Errorf("Subprocess finished with error: %v State: %v", err, cmd.ProcessState)
	}
}

func TestDefaultCommand_Start(t *testing.T) {
	// Set up mock execCommand and prep it to be returned
	execCommand = func(name string, args ...string) *exec.Cmd {
		return oldExecCommand("ls") // Every system has ls.
	}
	defer func() { execCommand = oldExecCommand }()

	b := &mock_bazel.MockBazel{}

	cmd := start(b, "//path/to:target", []string{"moo"})
	cmd.Start()

	if cmd.Stdout != os.Stdout {
		t.Errorf("Didn't set Stdout correctly")
	}
	if cmd.Stderr != os.Stderr {
		t.Errorf("Didn't set Stderr correctly")
	}
	if cmd.SysProcAttr.Setpgid != true {
		t.Errorf("Never set PGID (will prevent killing process trees -- see notes in ibazel.go")
	}

	b.AssertActions(t, [][]string{
		[]string{"Run", "--script_path=.*", "//path/to:target"},
	})
}
