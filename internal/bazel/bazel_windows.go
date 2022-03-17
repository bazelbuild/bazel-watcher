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
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

// Windows specific as it assures the bazelPath is always double quoted
// which assures we can support paths with whitespaces.
// It works by specifying the CmdLine after the exec command has been specified
// and by wrapping in double quotes the bazelPath content
//
// NOTE: SysProcAttr.CmdLine does not exist/is supported to be compiled on other
// OS other than Windows which is the reason why this new fn was created both for Windows and Unix
func setProcessAttributes(cmd *exec.Cmd, bazelPath string, args []string) {
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.CmdLine = fmt.Sprintf("%s %s", fmt.Sprintf("%q", bazelPath), strings.Join(args[:], " "))
}
