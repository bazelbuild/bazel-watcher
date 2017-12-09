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
	"os"
	"os/exec"
	"syscall"
)

type defaultCommand struct {
	target    string
	bazelArgs []string
	args      []string
	cmd       *exec.Cmd
}

// DefaultCommand is the normal mode of interacting with iBazel. If you start a
// server in this mode and notify of changes the server will be killed and
// restarted.
func DefaultCommand(bazelArgs []string, target string, args []string) Command {
	return &defaultCommand{
		target:    target,
		bazelArgs: bazelArgs,
		args:      args,
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

func (c *defaultCommand) Start() error {
	b := bazelNew()
	b.SetArguments(c.bazelArgs)

	b.WriteToStderr(true)
	b.WriteToStdout(true)

	c.cmd = start(b, c.target, c.args)

	c.cmd.Env = os.Environ()

	var err error
	if err = c.cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting process: %v\n", err)
		return err
	}
	fmt.Fprintf(os.Stderr, "Starting...")
	return nil
}

func (c *defaultCommand) NotifyOfChanges() {
	c.Terminate()
	c.Start()
}

func (c *defaultCommand) IsSubprocessRunning() bool {
	return subprocessRunning(c.cmd)
}
