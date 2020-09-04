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
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/bazelbuild/bazel-watcher/bazel"
	"github.com/bazelbuild/bazel-watcher/ibazel/log"
	"github.com/bazelbuild/bazel-watcher/ibazel/process_group"
)

var (
	execCommand  = process_group.Command
	bazelNew     = bazel.New
	waitDuration = flag.Duration(
		"graceful_termination_wait_duration",
		10*time.Second,
		"Specify the duration to wait for a graceful termination before sending SIGKILL to the subprocess")
)

// Command is an object that wraps the logic of running a task in Bazel and
// manipulating it.
type Command interface {
	Start() (*bytes.Buffer, error)
	Terminate()
	Kill()
	NotifyOfChanges() *bytes.Buffer
	IsSubprocessRunning() bool
}

// start will be called by most implementations since this logic is extremely
// common.
func start(b bazel.Bazel, target string, args []string) (*bytes.Buffer, process_group.ProcessGroup) {
	var filePattern strings.Builder
	filePattern.WriteString("bazel_script_path*")
	if runtime.GOOS == "windows" {
		filePattern.WriteString(".bat")
	}

	tmpfile, err := ioutil.TempFile("", filePattern.String())
	if err != nil {
		fmt.Print(err)
	}
	// Close the file so bazel can write over it
	if err := tmpfile.Close(); err != nil {
		fmt.Print(err)
	}

	// Start by building the binary
	_, outputBuffer, _ := b.Run("--script_path="+tmpfile.Name(), target)

	runScriptPath := tmpfile.Name()

	// Now that we have built the target, construct a executable form of it for
	// execution in a go routine.
	cmd := execCommand(runScriptPath, args...)
	cmd.RootProcess().Stdout = os.Stdout
	cmd.RootProcess().Stderr = os.Stderr

	return outputBuffer, cmd
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

func terminate(pg process_group.ProcessGroup) {
	pg.Signal(syscall.SIGTERM)
	done := make(chan bool, 1)
	go func() {
		select {
		case <-time.After(*waitDuration):
			log.Logf("The subprocess wasn't terminated within %s. Forcing to close.", *waitDuration)
			kill(pg)
		case <-done:
			// The subprocess was terminated with SIGTERM
		}
	}()
	pg.Wait()
	done <- true
	pg.Close()
}

func kill(pg process_group.ProcessGroup) {
	if subprocessRunning(pg.RootProcess()) {
		log.Logf("Sending SIGKILL to the subprocess")
		pg.Signal(syscall.SIGKILL)
	}
}
