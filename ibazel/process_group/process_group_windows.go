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

package process_group

import (
	"errors"
	"os/exec"
	"syscall"
)

const (
	createSuspended = 0x00000004
)

var (
	createJobObject          uintptr
	assignProcessToJobObject uintptr
	terminateJobObject       uintptr
	resumeThread             uintptr
)

type winProcessGroup struct {
	root *exec.Cmd
	job  uintptr
}

func init() {
	kernel32, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		panic(err)
	}

	createJobObject, err = syscall.GetProcAddress(kernel32, "CreateJobObjectW")
	if err != nil {
		panic(err)
	}

	assignProcessToJobObject, err = syscall.GetProcAddress(kernel32, "AssignProcessToJobObject")
	if err != nil {
		panic(err)
	}

	terminateJobObject, err = syscall.GetProcAddress(kernel32, "TerminateJobObject")
	if err != nil {
		panic(err)
	}

	resumeThread, err = syscall.GetProcAddress(kernel32, "ResumeThread")
	if err != nil {
		panic(err)
	}
}

// Command creates a new ProcessGroup with a root command specified by the
// arguments.
func Command(name string, arg ...string) ProcessGroup {
	root := exec.Command(name, arg...)
	root.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createSuspended}
	return &winProcessGroup{root, 0}
}

func (pg *winProcessGroup) RootProcess() *exec.Cmd {
	return pg.root
}

func (pg *winProcessGroup) Start() error {
	if pg.job != 0 {
		return errors.New("job already started")
	}

	err := pg.root.Start()
	if err != nil {
		return err
	}

	job, _, errno := syscall.Syscall(createJobObject, 2, 0, 0, 0)
	if errno != 0 {
		return errno
	} else if job == 0 {
		return errors.New("unknown error creating job")
	}

	ret, _, errno := syscall.Syscall(assignProcessToJobObject, 2, job, uintptr(pg.root.Process.Pid), 0)
	if errno != 0 {
		return errno
	} else if ret == 0 {
		return errors.New("unknown error assigning process to job")
	}

	ret, _, errno = syscall.Syscall(resumeThread, 1, uintptr(pg.root.Process.Pid), 0, 0)
	if errno != 0 {
		return errno
	} else if int(ret) < 0 {
		return errors.New("unknown error resuming process")
	}

	return nil
}

func (pg *winProcessGroup) Kill() error {
	if pg.job == 0 {
		return errors.New("job not started")
	}

	ret, _, errno := syscall.Syscall(terminateJobObject, 2, pg.job, 0, 0)
	if errno != 0 {
		return errno
	} else if ret == 0 {
		return errors.New("unknown error killing job")
	}

	return nil
}
