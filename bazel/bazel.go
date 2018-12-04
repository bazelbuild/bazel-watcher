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
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	blaze_query "github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf"
	"github.com/golang/protobuf/proto"
)

var BazelPath = flag.String("bazel_path", "bazel", "Path to the bazel binary to use for actions")

type Bazel interface {
	SetArguments([]string)
	WriteToStderr(v bool)
	WriteToStdout(v bool)
	Info() (map[string]string, error)
	Query(args ...string) (*blaze_query.QueryResult, error)
	Build(args ...string) (*bytes.Buffer, error)
	Test(args ...string) (*bytes.Buffer, error)
	Run(args ...string) (*exec.Cmd, *bytes.Buffer, error)
	Wait() error
	Cancel()
}

type bazel struct {
	cmd           *exec.Cmd
	args          []string
	ctx           context.Context
	cancel        context.CancelFunc
	writeToStderr bool
	writeToStdout bool
}

func New() Bazel {
	return &bazel{}
}

func (b *bazel) SetArguments(args []string) {
	b.args = args
}

// WriteToStderr when running an operation.
func (b *bazel) WriteToStderr(v bool) {
	b.writeToStderr = v
}

// WriteToStdout when running an operation.
func (b *bazel) WriteToStdout(v bool) {
	b.writeToStdout = v
}

func (b *bazel) newCommand(command string, args ...string) (*bytes.Buffer, *bytes.Buffer) {
	b.ctx, b.cancel = context.WithCancel(context.Background())

	args = append([]string{command}, args...)

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
	b.cmd = exec.CommandContext(b.ctx, *BazelPath, args...)

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

// Displays information about the state of the bazel process in the
// form of several "key: value" pairs.  This includes the locations of
// several output directories.  Because some of the
// values are affected by the options passed to 'bazel build', the
// info command accepts the same set of options.
//
// A single non-option argument may be specified (e.g. "bazel-bin"), in
// which case only the value for that key will be printed.
//
// The full list of keys and the meaning of their values is documented in
// the bazel User Manual, and can be programmatically obtained with
// 'bazel help info-keys'.
//
//   res, err := b.Info()
func (b *bazel) Info() (map[string]string, error) {
	b.WriteToStderr(false)
	b.WriteToStdout(false)
	stdoutBuffer, _ := b.newCommand("info")

	err := b.cmd.Run()
	if err != nil {
		return nil, err
	}
	return b.processInfo(stdoutBuffer.String())
}

func (b *bazel) processInfo(info string) (map[string]string, error) {
	lines := strings.Split(info, "\n")
	output := make(map[string]string, 0)
	for _, line := range lines {
		if line == "" || strings.Contains(line, "Starting local Bazel server and connecting to it...") {
			continue
		}
		data := strings.SplitN(line, ": ", 2)
		if len(data) < 2 {
			return nil, errors.New("Bazel info returned a non key-value pair")
		}
		output[data[0]] = data[1]
	}
	return output, nil
}

// Executes a query language expression over a specified subgraph of the
// build dependency graph.
//
// For example, to show all C++ test rules in the strings package, use:
//
//   res, err := b.Query('kind("cc_.*test", strings:*)')
//
// or to find all dependencies of //path/to/package:target, use:
//
//   res, err := b.Query('deps(//path/to/package:target)')
//
// or to find a dependency path between //path/to/package:target and //dependency:
//
//   res, err := b.Query('somepath(//path/to/package:target, //dependency)')
func (b *bazel) Query(args ...string) (*blaze_query.QueryResult, error) {
	blazeArgs := append([]string(nil), "--output=proto", "--order_output=no", "--color=no")
	blazeArgs = append(blazeArgs, args...)

	b.WriteToStderr(true)
	b.WriteToStdout(false)
	stdoutBuffer, _ := b.newCommand("query", blazeArgs...)

	err := b.cmd.Run()

	if err != nil {
		return nil, err
	}
	return b.processQuery(stdoutBuffer.Bytes())
}

func (b *bazel) processQuery(out []byte) (*blaze_query.QueryResult, error) {
	var qr blaze_query.QueryResult
	if err := proto.Unmarshal(out, &qr); err != nil {
		fmt.Fprintf(os.Stderr, "Could not read blaze query response. Error: %s\nOutput: %s\n", err, out)
		return nil, err
	}

	return &qr, nil
}

func (b *bazel) Build(args ...string) (*bytes.Buffer, error) {
	stdoutBuffer, stderrBuffer := b.newCommand("build", append(b.args, args...)...)
	err := b.cmd.Run()

	_, _ = stdoutBuffer.Write(stderrBuffer.Bytes())
	return stdoutBuffer, err
}

func (b *bazel) Test(args ...string) (*bytes.Buffer, error) {
	stdoutBuffer, stderrBuffer := b.newCommand("test", append(b.args, args...)...)
	err := b.cmd.Run()

	_, _ = stdoutBuffer.Write(stderrBuffer.Bytes())
	return stdoutBuffer, err
}

// Build the specified target (singular) and run it with the given arguments.
func (b *bazel) Run(args ...string) (*exec.Cmd, *bytes.Buffer, error) {
	b.WriteToStderr(true)
	b.WriteToStdout(true)
	stdoutBuffer, stderrBuffer := b.newCommand("run", args...)
	b.cmd.Stdin = os.Stdin

	_, _ = stdoutBuffer.Write(stderrBuffer.Bytes())

	err := b.cmd.Run()
	if err != nil {
		return nil, stdoutBuffer, err
	}

	return b.cmd, stdoutBuffer, err
}

func (b *bazel) Wait() error {
	res := b.cmd.Wait()
	if res.Error() == "exec: Wait was already called" {
		if b.cmd.ProcessState.Success() {
			return nil
		}
	}
	return res
}

// Cancel the currently running operation. Useful if you call Run(target) and
// would like to stop the action running in a goroutine.
func (b *bazel) Cancel() {
	if b.cancel == nil {
		return
	}

	b.cancel()
}
