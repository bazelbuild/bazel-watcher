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
	"fmt"
	"io"
	"os"

	"github.com/bazelbuild/bazel-watcher/ibazel/process_group"
)

type notifyCommand struct {
	target    string
	bazelArgs []string
	args      []string

	pg    process_group.ProcessGroup
	stdin io.WriteCloser
}

// NotifyCommand is an alternate mode for starting a command. In this mode the
// command will be notified on stdin that the source files have changed.
func NotifyCommand(bazelArgs []string, target string, args []string) Command {
	return &notifyCommand{
		target:    target,
		bazelArgs: bazelArgs,
		args:      args,
	}
}

func (c *notifyCommand) Terminate() {
	if c.pg != nil && !subprocessRunning(c.pg.RootProcess()) {
		return
	}

	// Kill it with fire by sending SIGKILL to the process PID which should
	// propagate down to any subprocesses in the PGID (Process Group ID). To
	// send to the PGID, send the signal to the negative of the process PID.
	// Normally I would do this by calling c.cmd.Process.Signal, but that
	// only goes to the PID not the PGID.
	c.pg.Kill()
	c.pg.RootProcess().Wait()
	c.pg = nil
}

func (c *notifyCommand) Start() (*bytes.Buffer, error) {
	b := bazelNew()
	b.SetArguments(c.bazelArgs)

	b.WriteToStderr(true)
	b.WriteToStdout(true)

	var outputBuffer *bytes.Buffer
	outputBuffer, c.pg = start(b, c.target, c.args)
	// Keep the writer around.
	var err error
	c.stdin, err = c.pg.RootProcess().StdinPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting stdin pipe: %v\n", err)
		return outputBuffer, err
	}

	c.pg.RootProcess().Env = append(os.Environ(), "IBAZEL_NOTIFY_CHANGES=y")

	if err = c.pg.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting process: %v\n", err)
		return outputBuffer, err
	}
	fmt.Fprintf(os.Stderr, "Starting...")
	return outputBuffer, nil
}

func (c *notifyCommand) NotifyOfChanges() *bytes.Buffer {
	b := bazelNew()
	b.SetArguments(c.bazelArgs)

	b.WriteToStderr(true)
	b.WriteToStdout(true)

	outputBuffer, res := b.Build(c.target)
	if res != nil {
		fmt.Fprintf(os.Stderr, "FAILURE: %v\n", res)
		_, err := c.stdin.Write([]byte("IBAZEL_BUILD_COMPLETED FAILURE\n"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing failure to stdin: %s\n%v", err)
		}
	} else {
		fmt.Fprintf(os.Stderr, "SUCCESS\n")
		_, err := c.stdin.Write([]byte("IBAZEL_BUILD_COMPLETED SUCCESS\n"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing success to stdin: %v\n", err)
		}
	}
	return outputBuffer
}

func (c *notifyCommand) IsSubprocessRunning() bool {
	return c.pg != nil && subprocessRunning(c.pg.RootProcess())
}
