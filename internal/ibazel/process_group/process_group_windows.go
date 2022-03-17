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
	"os/exec"
	"syscall"
	"unsafe"
)

type winProcessGroup struct {
	root   *exec.Cmd
	job    syscall.Handle
	ioport syscall.Handle
}

// Command creates a new ProcessGroup with a root command specified by the
// arguments.
func Command(name string, arg ...string) ProcessGroup {
	root := exec.Command(name, arg...)
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

	pg.job, err = createJobObject()
	if err != nil {
		return err
	}

	pg.ioport, err = syscall.CreateIoCompletionPort(syscall.InvalidHandle, syscall.Handle(0), 0, 1)
	if err != nil {
		return err
	}

	port := jobObjectAssociationCompletionPort{
		CompletionKey:  pg.job,
		CompletionPort: pg.ioport,
	}

	err = setInformationJobObject(pg.job, jobObjectAssociateCompletionPortInformation, uintptr(unsafe.Pointer(&port)), unsafe.Sizeof(port))
	if err != nil {
		return err
	}

	process, err := syscall.OpenProcess(processAllAccess, false, uint32(pg.root.Process.Pid))
	if err != nil {
		return err
	}

	err = assignProcessToJobObject(pg.job, process)
	if err != nil {
		return err
	}

	err = ntResumeProcess(process)
	if err != nil {
		return err
	}

	syscall.CloseHandle(process)

	return nil
}

func (pg *winProcessGroup) Signal(signum syscall.Signal) error {
	// signum is ignored on Windows as there's no support for signals
	if pg.job == 0 {
		return errors.New("job not started")
	}

	err := terminateJobObject(pg.job, 0)
	if err != nil {
		return err
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

func (pg *winProcessGroup) Close() error {
	err := syscall.CloseHandle(pg.job)
	if err != nil {
		return err
	}

	err = syscall.CloseHandle(pg.ioport)
	if err != nil {
		return err
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
