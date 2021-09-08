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

// +build !windows

package bazel

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func (b *bazel) newCommand(command string, args ...string) (*bytes.Buffer, *bytes.Buffer) {
	b.ctx, b.cancel = context.WithCancel(context.Background())

	args = append([]string{command}, args...)
	args = append(b.startupArgs, args...)

	if b.writeToStderr || b.writeToStdout {
		containsColor := false
		for _, arg := range args {
			if strings.HasPrefix(arg, "--color") {
				containsColor = true
			}
		}
		if !containsColor {
			args = append(args, "--color=yes")
		}
	}

	b.cmd = exec.CommandContext(b.ctx, findBazel(), args...)
	// Always set a process group for unix based systems
	b.cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		Pgid:    0,
	}

	stdoutBuffer := new(bytes.Buffer)
	stderrBuffer := new(bytes.Buffer)
	if b.writeToStdout {
		b.cmd.Stdout = io.MultiWriter(os.Stdout, stdoutBuffer)
	} else {
		b.cmd.Stdout = stdoutBuffer
	}
	if b.writeToStderr {
		b.cmd.Stderr = io.MultiWriter(os.Stderr, stderrBuffer)
	} else {
		b.cmd.Stderr = stderrBuffer
	}

	return stdoutBuffer, stderrBuffer
}

func (b *bazel) BuildCancelable(cancelCh chan bool, args ...string) (*bytes.Buffer, error) {
	stdoutBuffer, stderrBuffer := b.newCommand("build", append(b.args, args...)...)
	doneCh := make(chan error)

	go func() {
		doneCh <- b.cmd.Run()
	}()

	select {
	case e := <-doneCh:
		_, _ = stdoutBuffer.Write(stderrBuffer.Bytes())

		return stdoutBuffer, e
	case <-cancelCh:
		b.cmd.Process.Signal(syscall.SIGTERM)
		<-doneCh

		return nil, nil
	}
}

func (b *bazel) TestCancelable(cancelCh chan bool, args ...string) (*bytes.Buffer, error) {
	stdoutBuffer, stderrBuffer := b.newCommand("test", append(b.args, args...)...)
	doneCh := make(chan error)

	go func() {
		doneCh <- b.cmd.Run()
	}()

	select {
	case e := <-doneCh:
		_, _ = stdoutBuffer.Write(stderrBuffer.Bytes())

		return stdoutBuffer, e
	case <-cancelCh:
		b.cmd.Process.Signal(syscall.SIGTERM)
		<-doneCh

		return nil, nil
	}
}
