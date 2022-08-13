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
	"os/exec"
	"testing"

	"github.com/bazelbuild/bazel-watcher/internal/bazel"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/process_group"
)

var oldExecCommand = execCommand
var oldBazelNew = bazel.New

func assertKilled(t *testing.T, cmd *exec.Cmd) {
	t.Helper()
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
	execCommand = func(name string, args ...string) process_group.ProcessGroup {
		return oldExecCommand("ls") // Every system has ls.
	}
	defer func() { execCommand = oldExecCommand }()

	if subprocessRunning(nil) {
		t.Errorf("Nil subprocesses don't run")
	}

	cmd := exec.Command("sleep", ".1")

	if subprocessRunning(cmd) {
		t.Errorf("New subprocess shouldn't have been started yet. State: %v", cmd.ProcessState)
	}

	if err := cmd.Start(); err != nil {
		t.Errorf("cmd.Start(): %v", err)
	}

	if !subprocessRunning(cmd) {
		t.Errorf("New subprocess was never started. State: %v", cmd.ProcessState)
	}

	err := cmd.Wait()
	if err != nil {
		t.Errorf("Subprocess finished with error: %v State: %v", err, cmd.ProcessState)
	} else if subprocessRunning(cmd) {
		t.Errorf("Subprocess still running State: %v", cmd.ProcessState)
	}
}
