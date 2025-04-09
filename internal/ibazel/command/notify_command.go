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
	"io"
	"os"
	"sync"

	"github.com/bazelbuild/bazel-watcher/internal/ibazel/log"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/process_group"
)

type notifyCommand struct {
	target      string
	startupArgs []string
	bazelArgs   []string
	args        []string
	pg          process_group.ProcessGroup
	stdin       io.WriteCloser
	termSync    sync.Once
}

// NotifyCommand is an alternate mode for starting a command. In this mode the
// command will be notified on stdin that the source files have changed.
func NotifyCommand(startupArgs []string, bazelArgs []string, target string, args []string) Command {
	return &notifyCommand{
		startupArgs: startupArgs,
		target:      target,
		bazelArgs:   bazelArgs,
		args:        args,
	}
}

func (c *notifyCommand) Terminate() {
	if !c.IsSubprocessRunning() {
		c.pg = nil
		return
	}
	c.termSync.Do(func() {
		terminate(c.pg)
	})
	c.pg = nil
}

func (c *notifyCommand) Kill() {
	if c.pg != nil {
		kill(c.pg)
	}
}

func (c *notifyCommand) Start() (*bytes.Buffer, error) {
	b := bazelNew()
	b.SetStartupArgs(c.startupArgs)
	b.SetArguments(c.bazelArgs)

	b.WriteToStderr(true)
	b.WriteToStdout(true)

	var outputBuffer *bytes.Buffer
	outputBuffer, c.pg = start(b, c.target, c.args)
	// Keep the writer around.
	var err error
	c.stdin, err = c.pg.RootProcess().StdinPipe()
	if err != nil {
		log.Errorf("Error getting stdin pipe: %v", err)
		return outputBuffer, err
	}

	c.pg.RootProcess().Env = append(os.Environ(), "IBAZEL_NOTIFY_CHANGES=y")

	if err = c.pg.Start(); err != nil {
		log.Errorf("Error starting process: %v", err)
		return outputBuffer, err
	}
	log.Log("Starting...")
	c.termSync = sync.Once{}
	return outputBuffer, nil
}

func (c *notifyCommand) NotifyOfChanges() *bytes.Buffer {
	b := bazelNew()
	b.SetStartupArgs(c.startupArgs)
	// We must always pass --remote_download_outputs=toplevel to ensure that the runfiles of the target are correctly downloaded before we notify the process
	b.SetArguments(append(c.bazelArgs, "--remote_download_outputs=toplevel"))

	b.WriteToStderr(true)
	b.WriteToStdout(true)

	_, err := c.stdin.Write([]byte("IBAZEL_BUILD_STARTED\n"))
	if err != nil {
		log.Errorf("Error writing build to stdin: %s", err)
	}

	outputBuffer, res := b.Build(c.target)
	if res != nil {
		log.Errorf("IBAZEL BUILD FAILURE: %v", res)
		_, err := c.stdin.Write([]byte("IBAZEL_BUILD_COMPLETED FAILURE\n"))
		if err != nil {
			log.Errorf("Error writing failure to stdin: %s", err)
		}
	} else {
		log.Log("IBAZEL BUILD SUCCESS")
		_, err := c.stdin.Write([]byte("IBAZEL_BUILD_COMPLETED SUCCESS\n"))
		if err != nil {
			log.Errorf("Error writing success to stdin: %v", err)
		}
		if !c.IsSubprocessRunning() {
			log.Log("Restarting process...")
			c.Terminate()
			c.Start()
		}
	}
	return outputBuffer
}

func (c *notifyCommand) IsSubprocessRunning() bool {
	return c.pg != nil && subprocessRunning(c.pg.RootProcess())
}
