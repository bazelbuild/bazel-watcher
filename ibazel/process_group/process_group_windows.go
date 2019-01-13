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
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"syscall"
	"unsafe"
)

const (
	createSuspended     = 0x00000004
	threadSuspendResume = 0x0002
	processAllAcccess   = 0x1F0FFF

	jobObjectAssociateCompletionPortInformation = 7

	jobObjectMsgActiveProcessZero = 4
)

var (
	createJobObject          uintptr
	setInformationJobObject  uintptr
	assignProcessToJobObject uintptr
	terminateJobObject       uintptr
	ntResumeProcess          uintptr
)

type winProcessGroup struct {
	root   *exec.Cmd
	job    syscall.Handle
	ioport syscall.Handle
}

type threadEntry32 struct {
	dwSize             uint32
	cntUsage           uint32
	th32ThreadID       uint32
	th32OwnerProcessID uint32
	tpBasePri          uint32
	tpDeltaPri         uint32
	dwFlags            uint32
}

type jobObjectAssociationCompletionPort struct {
	CompletionKey  uintptr
	CompletionPort syscall.Handle
}

func init() {
	kernel32, err := syscall.LoadLibrary("kernel32.dll")
	if err != nil {
		panic(err)
	}

	ntdll, err := syscall.LoadLibrary("ntdll.dll")
	if err != nil {
		panic(err)
	}

	createJobObject, err = syscall.GetProcAddress(kernel32, "CreateJobObjectW")
	if err != nil {
		panic(err)
	}

	setInformationJobObject, err = syscall.GetProcAddress(kernel32, "SetInformationJobObject")
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

	ntResumeProcess, err = syscall.GetProcAddress(ntdll, "NtResumeProcess")
	if err != nil {
		panic(err)
	}
}

// Command creates a new ProcessGroup with a root command specified by the
// arguments.
func Command(name string, arg ...string) ProcessGroup {
	root := exec.Command(name, arg...)
	fmt.Println(name, arg)
	root.SysProcAttr = &syscall.SysProcAttr{CreationFlags: createSuspended}
	return &winProcessGroup{root, syscall.Handle(0), syscall.Handle(0)}
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
	}
	pg.job = syscall.Handle(job)

	pg.ioport, err = syscall.CreateIoCompletionPort(syscall.InvalidHandle, syscall.Handle(0), 0, 1)
	if err != nil {
		return err
	}

	port := jobObjectAssociationCompletionPort{
		CompletionKey:  job,
		CompletionPort: pg.ioport,
	}

	_, _, errno = syscall.Syscall6(setInformationJobObject, 4, uintptr(pg.job), jobObjectAssociateCompletionPortInformation, uintptr(unsafe.Pointer(&port)), unsafe.Sizeof(port), 0, 0)
	if errno != 0 {
		return errno
	}

	phandle, err := syscall.OpenProcess(processAllAcccess, false, uint32(pg.root.Process.Pid))
	if err != nil {
		return err
	}

	_, _, errno = syscall.Syscall(assignProcessToJobObject, 2, uintptr(pg.job), uintptr(phandle), 0)
	if errno != 0 {
		return errno
	}

	_, _, errno = syscall.Syscall(ntResumeProcess, 1, uintptr(phandle), 0, 0)
	if errno != 0 {
		return errno
	}

	return nil
}

func (pg *winProcessGroup) Kill() error {
	if pg.job == 0 {
		return errors.New("job not started")
	}

	ret, _, errno := syscall.Syscall(terminateJobObject, 2, uintptr(pg.job), 0, 0)
	if errno != 0 {
		return errno
	} else if ret == 0 {
		return errors.New("unknown error killing job")
	}

	return nil
}

func (pg *winProcessGroup) Wait() error {
	var code uint32
	var key uint32
	var op *syscall.Overlapped
	for {
		err := syscall.GetQueuedCompletionStatus(pg.ioport, &code, &key, &op, syscall.INFINITE)
		if err != nil {
			return err
		}
		if key == uint32(pg.job) && code == jobObjectMsgActiveProcessZero {
			break
		}
	}

	return nil
}

func (pg *winProcessGroup) Run() error {
	if err := pg.Start(); err != nil {
		return err
	}
	return pg.Wait()
}

func (pg *winProcessGroup) CombinedOutput() ([]byte, error) {
	var b bytes.Buffer
	pg.root.Stdout = &b
	pg.root.Stderr = &b
	err := pg.Run()
	return b.Bytes(), err
}
