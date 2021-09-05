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

package bazel

import (
	"bytes"
	"context"
	"fmt"
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

	bazelPath := findBazel()
	b.cmd = exec.CommandContext(b.ctx, bazelPath, args...)

	// windows specific as it assures the bazelPath is always double quoted
	// which assures we can support paths with whitespaces.
	// It works by specifying the CmdLine after the exec command has been specified
	// and by wrapping in double quotes the bazelPath content
	//
	// NOTE: SysProcAttr.CmdLine does not exist/is supported to be compiled on other
	// OS other than Windows which is the reason why newCommand fn was created both
	// for Windows and Unix
	bazelPath = fmt.Sprintf("\"%s\"", bazelPath)
	b.cmd.SysProcAttr = &syscall.SysProcAttr{}
	b.cmd.SysProcAttr.CmdLine = fmt.Sprintf("%s %s", bazelPath, strings.Join(args[:], " "))

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
