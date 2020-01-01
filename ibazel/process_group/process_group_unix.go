// Copyright 2018 The Bazel Authors. All rights reserved.
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

package process_group

import (
	"os/exec"
	"syscall"
)

type unixProcessGroup struct {
	root *exec.Cmd
}

// Command creates a new ProcessGroup with a root command specified by the
// arguments.
func Command(name string, arg ...string) ProcessGroup {
	root := exec.Command(name, arg...)
	root.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
		// TODO: Consider supporting Pdeathsig for Linux.
		// Consequence of not having Pdeathdig is that on non-Linux systems,
		// if SIGTERM doesn't stop child procs then they may "leak" and be
		// reparented 'up the chain' somewhere when the minion process
		// terminates.
		// Pdeathsig: syscall.SIGKILL
	}
	return &unixProcessGroup{root}
}

func (pg *unixProcessGroup) RootProcess() *exec.Cmd {
	return pg.root
}

func (pg *unixProcessGroup) Start() error {
	return pg.root.Start()
}

func (pg *unixProcessGroup) Kill() error {
	// Kill all proccesses in group.
	pgid, err := syscall.Getpgid(pg.root.Process.Pid)
	if err == nil {
		return syscall.Kill(-pgid, syscall.SIGKILL)
	}

	// Fallback to assuming the Pid is the process group ID.
	return syscall.Kill(-pg.root.Process.Pid, syscall.SIGKILL)
}

func (pg *unixProcessGroup) Wait() error {
	return pg.root.Wait()
}

func (pg *unixProcessGroup) Close() error {
	return nil
}

func (pg *unixProcessGroup) CombinedOutput() ([]byte, error) {
	return pg.root.CombinedOutput()
}
