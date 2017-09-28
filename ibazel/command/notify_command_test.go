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
	"os/exec"
	"syscall"
	"testing"

	"github.com/bazelbuild/bazel-watcher/bazel"
	mock_bazel "github.com/bazelbuild/bazel-watcher/bazel/testing"
)

func TestNotifyCommand(t *testing.T) {
	cmd := exec.Command("cat")
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	c := &notifyCommand{
		args:      []string{"moo"},
		bazelArgs: []string{},
		cmd:       cmd,
		target:    "//path/to:target",
	}

	if c.IsSubprocessRunning() {
		t.Errorf("New subprocess shouldn't have been started yet. State: %v", cmd.ProcessState)
	}

	var err error
	c.stdin, err = cmd.StdinPipe()
	if err != nil {
		t.Error(err)
	}

	// Mock out bazel to return non-error on test
	b := &mock_bazel.MockBazel{}
	b.WaitError(nil)
	bazelNew = func() bazel.Bazel { return b }
	defer func() { bazelNew = oldBazelNew }()

	c.NotifyOfChanges()
	b.WaitError(errors.New("Demo error"))
	c.NotifyOfChanges()
	b.WaitError(nil)
	c.NotifyOfChanges()

	b.AssertActions(t, [][]string{
		[]string{"WriteToStderr"},
		[]string{"WriteToStdout"},
		[]string{"Build", "//path/to:target"},
		[]string{"WriteToStderr"},
		[]string{"WriteToStdout"},
		[]string{"Build", "//path/to:target"},
		[]string{"WriteToStderr"},
		[]string{"WriteToStdout"},
		[]string{"Build", "//path/to:target"},
	})

	err = c.stdin.Close()
	if err != nil {
		t.Error(err)
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Error(err)
	}

	// Read on the pipe is only valid in between start and wait so read now.
	expected := "IBAZEL_BUILD_COMPLETED SUCCESS\nIBAZEL_BUILD_COMPLETED FAILURE\nIBAZEL_BUILD_COMPLETED SUCCESS\n"
	if expected != string(out) {
		t.Errorf("Not equal.\nGot:  %s\nWant: %s", string(out), expected)
	}
}
