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
	"bytes"
	"os"
	"os/exec"
	"syscall"

	"github.com/bazelbuild/bazel-watcher/bazel"
)

var execCommand = exec.Command
var bazelNew = bazel.New

// Command is an object that wraps the logic of running a task in Bazel and
// manipulating it.
type Command interface {
	Start() (*bytes.Buffer, error)
	Terminate()
	NotifyOfChanges() *bytes.Buffer
	IsSubprocessRunning() bool
}

// start will be called by most implementations since this logic is extremely
// common.
func start(b bazel.Bazel, target string, args []string) (*bytes.Buffer, *exec.Cmd) {
	// Build and run the target in a go routine with bazel. Since the direct_run
	// functionaliy was made default that works fine.
	args = append([]string{"run", target}, args...)
	cmd := execCommand(*bazel.BazelPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set a process group id (PGID) on the subprocess. This is
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return nil, cmd
}

func subprocessRunning(cmd *exec.Cmd) bool {
	if cmd == nil {
		return false
	}
	if cmd.Process == nil {
		return false
	}
	if cmd.ProcessState != nil {
		if cmd.ProcessState.Exited() {
			return false
		}
	}

	return true
}
