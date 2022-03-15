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
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	
	"github.com/bazelbuild/bazel-watcher/ibazel/log"
)

// Windows specific as it assures the bazelPath is always double quoted
// which assures we can support paths with whitespaces.
// It works by specifying the CmdLine after the exec command has been specified
// and by wrapping in double quotes the bazelPath content
//
// NOTE: SysProcAttr.CmdLine does not exist/is supported to be compiled on other
// OS other than Windows which is the reason why this new fn was created both for Windows and Unix
func setProcessAttributes(cmd *exec.Cmd, bazelPath string, args []string) {
	// Flag to always run this cmd as a new windows process group
	// otherwise we cannot correctly send the ctrl+c event on
	// canceleable methods
	cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
	cmd.SysProcAttr.CmdLine = fmt.Sprintf("%s %s", fmt.Sprintf("%q", bazelPath), strings.Join(args[:], " "))
}

// Internal function useful to send CTRL+C events
// to a given process on Windows OS
func sendCtrlBreak(pid int) {
	kernel32DLL, dllLoadingError := syscall.LoadDLL("kernel32.dll")
	if dllLoadingError != nil {
		log.Errorf("Error loading a required dll `kernel32.dll` when `IBAZEL_ABORT_COMPILATION_EARLY=1`: %v\n", dllLoadingError)
	}
	GenerateConsoleCtrlEventProc, findProcError := kernel32DLL.FindProc("GenerateConsoleCtrlEvent")
	if findProcError != nil {
		log.Errorf("Error calling FindProc with GenerateConsoleCtrlEvent: %v\n", findProcError)
	}
	ConsoleCtrlEventResult, _, eventError := GenerateConsoleCtrlEventProc.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
	if ConsoleCtrlEventResult == 0 {
		log.Errorf("Error while generating a ConsoleCtrlEvent on pid %d: %v\n", pid, eventError)
	}
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
		sendCtrlBreak(b.cmd.Process.Pid)
		b.Cancel()
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
		sendCtrlBreak(b.cmd.Process.Pid)
		b.Cancel()
		<-doneCh

		return nil, nil
	}
}
