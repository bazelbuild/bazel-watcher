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

// Implements a platform-independent process group. In effect, this allows you
// to terminate an entire tree of processes in one go. On Linux, this uses the
// Process Group system. On Windows, this uses Job Objects.
//
// Most of the things you would normally do with an exec.Cmd are safe to do
// with the RootProcess() Cmd, with two exceptions:
//
// - You cannot call .Start() on it. Use ProcessGroup.Start() instead. Doing
//   so will work on Linux but not Windows.
// - You should not change .SysProcAttr.

package process_group

import (
	"os/exec"
	"syscall"
)

// ProcessGroup represents a tree of processes that can be terminated
// simultaneously.
type ProcessGroup interface {
	RootProcess() *exec.Cmd
	Start() error
	Signal(signum syscall.Signal) error
	Wait() error
	Close() error
	CombinedOutput() ([]byte, error)
}
