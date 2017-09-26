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
	"syscall"

	"github.com/bazelbuild/bazel-watcher/bazel"
)

type defaultCommand struct {
	target string
	b      bazel.Bazel
	args   []string
	cmd    *exec.Cmd
}

// DefaultCommand is the normal mode of interacting with iBazel. If you start a
// server in this mode and notify of changes the server will be killed and
// restarted.
func DefaultCommand(bazel bazel.Bazel, target string, args []string) Command {
	return &defaultCommand{
		target: target,
		b:      bazel,
		args:   args,
	}
}

func (c *defaultCommand) Terminate() {
	if !subprocessRunning(c.cmd) {
		return
	}

	// Kill it with fire by sending SIGKILL to the process PID which should
	// propagate down to any subprocesses in the PGID (Process Group ID). To
	// send to the PGID, send the signal to the negative of the process PID.
	// Normally I would do this by calling c.cmd.Process.Signal, but that
	// only goes to the PID not the PGID.
	syscall.Kill(-c.cmd.Process.Pid, syscall.SIGKILL)
	c.cmd.Wait()
	c.cmd = nil
}

func (c *defaultCommand) Start() {
	c.cmd = start(c.b, c.target, c.args)
}

func (c *defaultCommand) NotifyOfChanges() {
	c.Terminate()
}

func (c *defaultCommand) IsSubprocessRunning() bool {
	return subprocessRunning(c.cmd)
}
