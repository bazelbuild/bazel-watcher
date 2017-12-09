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
	"fmt"
	"io/ioutil"
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
	Start() error
	Terminate()
	NotifyOfChanges()
	IsSubprocessRunning() bool
}

// start will be called by most implementations since this logic is extremely
// common.
func start(b bazel.Bazel, target string, args []string) *exec.Cmd {
	tmpfile, err := ioutil.TempFile("", "bazel_script_path")
	if err != nil {
		fmt.Print(err)
	}
	// Close the file so bazel can write over it
	if err := tmpfile.Close(); err != nil {
		fmt.Print(err)
	}

	// Start by building the binary
	b.Run("--script_path="+tmpfile.Name(), target)

	runScriptPath := tmpfile.Name()

	// Now that we have built the target, construct a executable form of it for
	// execution in a go routine.
	cmd := execCommand(runScriptPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Set a process group id (PGID) on the subprocess. This is
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	return cmd
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
