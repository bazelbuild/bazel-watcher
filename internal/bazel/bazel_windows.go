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
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func setProcessAttributes(cmd *exec.Cmd, bazelPath string, args []string) {
	// Let os/exec quote native executable arguments using Windows argv rules.
	extension := strings.ToLower(filepath.Ext(cmd.Path))
	if cmd.Err != nil || (extension != ".cmd" && extension != ".bat") {
		return
	}
	for _, arg := range append([]string{cmd.Path}, args...) {
		if strings.ContainsAny(arg, "\r\n") {
			cmd.Err = fmt.Errorf("Windows batch wrappers cannot accept arguments containing newlines")
			return
		}
	}
	command := escapeBatchMetacharacters(cmd.Path)
	for _, arg := range args {
		command += " " + escapeBatchArgument(arg)
	}
	cmd.Path = filepath.Join(os.Getenv("SystemRoot"), "System32", "cmd.exe")
	cmd.Args = []string{cmd.Path, "/d", "/s", "/v:off", "/c", command}
	cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: `/d /s /v:off /c "` + command + `"`}
	// Killing only cmd.exe leaves the wrapped Bazel client running.
	cmd.Cancel = func() error {
		killer := filepath.Join(os.Getenv("SystemRoot"), "System32", "taskkill.exe")
		if err := exec.Command(killer, "/t", "/f", "/pid", strconv.Itoa(cmd.Process.Pid)).Run(); err != nil {
			return cmd.Process.Kill()
		}
		return nil
	}
}

func escapeBatchArgument(arg string) string {
	// Quote for the child argv parser, then protect both cmd.exe parsing passes:
	// the /c invocation and the batch wrapper's forwarding of %*.
	var quoted strings.Builder
	quoted.WriteByte('"')
	slashes := 0
	for _, char := range arg {
		if char == '\\' {
			slashes++
			continue
		}
		if char == '"' {
			quoted.WriteString(strings.Repeat("\\", slashes*2+1))
		} else {
			quoted.WriteString(strings.Repeat("\\", slashes))
		}
		slashes = 0
		quoted.WriteRune(char)
	}
	quoted.WriteString(strings.Repeat("\\", slashes*2))
	quoted.WriteByte('"')
	return escapeBatchMetacharacters(escapeBatchMetacharacters(quoted.String()))
}

func escapeBatchMetacharacters(arg string) string {
	var result strings.Builder
	for _, char := range arg {
		if strings.ContainsRune("()[]%!^\"`<>&|;, *?", char) {
			result.WriteByte('^')
		}
		result.WriteRune(char)
	}
	return result.String()
}
