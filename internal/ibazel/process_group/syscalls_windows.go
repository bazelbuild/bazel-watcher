package process_group

import (
	"syscall"
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")
	modntdll    = syscall.NewLazyDLL("ntdll.dll")

	procCreateJobObject          = modkernel32.NewProc("CreateJobObjectW")
	procSetInformationJobObject  = modkernel32.NewProc("SetInformationJobObject")
	procAssignProcessToJobObject = modkernel32.NewProc("AssignProcessToJobObject")
	procTerminateJobObject       = modkernel32.NewProc("TerminateJobObject")
	procNtResumeProcess          = modntdll.NewProc("NtResumeProcess")
)

const (
	// CREATE_SUSPEND flag - used in CreateProcess 'dwCreationFlags' field to
	// specify that the process should start suspended.
	createSuspended = 0x00000004

	// PROCESS_ALL_ACCESS flag - used in OpenProcess 'dwDesiredAccess' field to
	// request all permission. Seems needed for NtResumeProcess.
	processAllAccess = 0x1F0FFF

	// JobObjectAssociateCompletionPortInformation - used in
	// SetInformationJobObject to set completion port of job object.
	jobObjectAssociateCompletionPortInformation = 7

	// JOB_OBJECT_MSG_ACTIVE_PROCESS_ZERO - message returned from IO completion
	// port when no more processes are active in job object.
	jobObjectMsgActiveProcessZero = 4
)

type jobObjectAssociationCompletionPort struct {
	CompletionKey  syscall.Handle
	CompletionPort syscall.Handle
}

func createJobObject() (syscall.Handle, error) {
	job, _, errno := syscall.Syscall(procCreateJobObject.Addr(), 2, 0, 0, 0)
	if errno != 0 {
		return 0, errno
	}
	return syscall.Handle(job), nil
}

func setInformationJobObject(job syscall.Handle, infoClass int, objInfo uintptr, objInfoLen uintptr) error {
	_, _, errno := syscall.Syscall6(procSetInformationJobObject.Addr(), 4, uintptr(job), uintptr(infoClass), objInfo, objInfoLen, 0, 0)
	if errno != 0 {
		return errno
	}
	return nil
}

func assignProcessToJobObject(job syscall.Handle, process syscall.Handle) error {
	_, _, errno := syscall.Syscall(procAssignProcessToJobObject.Addr(), 2, uintptr(job), uintptr(process), 0)
	if errno != 0 {
		return errno
	}
	return nil
}

func terminateJobObject(job syscall.Handle, exitCode uint) error {
	_, _, errno := syscall.Syscall(procTerminateJobObject.Addr(), 2, uintptr(job), uintptr(exitCode), 0)
	if errno != 0 {
		return errno
	}
	return nil
}

func ntResumeProcess(process syscall.Handle) error {
	_, _, errno := syscall.Syscall(procNtResumeProcess.Addr(), 1, uintptr(process), 0, 0)
	if errno != 0 {
		return errno
	}
	return nil
}
