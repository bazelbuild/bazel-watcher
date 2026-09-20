// Copyright 2026 The Bazel Authors. All rights reserved.
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
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"syscall"
	"testing"
	"time"
)

func TestWindowsBazelHelper(t *testing.T) {
	if os.Getenv("IBAZEL_ARGV_HELPER") != "1" {
		return
	}
	for i, arg := range os.Args {
		if arg == "--" {
			if os.Getenv("IBAZEL_HELPER_SLEEP") == "1" {
				json.NewEncoder(os.Stdout).Encode(os.Getpid())
				time.Sleep(time.Minute)
				os.Exit(24)
			}
			json.NewEncoder(os.Stdout).Encode(os.Args[i+1:])
			os.Exit(23)
		}
	}
	os.Exit(24)
}

func TestWindowsBazelArguments(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(t.TempDir(), "wrappers & %IBAZEL_EXPAND_ME% !")
	if err := os.Mkdir(dir, 0700); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(executable)
	if err != nil {
		t.Fatal(err)
	}
	helper := filepath.Join(dir, "native helper.exe")
	if err := os.WriteFile(helper, contents, 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("IBAZEL_ARGV_HELPER", "1")
	t.Setenv("IBAZEL_TEST_EXE", helper)
	t.Setenv("IBAZEL_EXPAND_ME", "unexpected expansion")
	want := []string{"", "plain", "two words", `kind(".* rule", //...)`, `C:\path with spaces\`, `a\"b`, "a&b|c<d>e", `quote"&echo injected&"`, "//pkg:target", "(x)^y!z", "%IBAZEL_EXPAND_ME%", "日本語"}
	for _, extension := range []string{".exe", ".cmd", ".BAT"} {
		t.Run(extension, func(t *testing.T) {
			path := helper
			if extension != ".exe" {
				path = filepath.Join(dir, "bazel wrapper"+extension)
				if err := os.WriteFile(path, []byte("@\"%IBAZEL_TEST_EXE%\" %*\r\n@exit /b %errorlevel%\r\n"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			old := *bazelPathFlag
			*bazelPathFlag = path
			t.Cleanup(func() { *bazelPathFlag = old })
			b := &bazel{}
			stdout, stderr := b.newCommand("-test.run=TestWindowsBazelHelper", append([]string{"--"}, want...)...)
			defer b.cancel()
			err := b.cmd.Run()
			var exit *exec.ExitError
			if !errors.As(err, &exit) || exit.ExitCode() != 23 {
				t.Fatalf("exit: %v; stderr: %s; stdout: %s", err, stderr, stdout)
			}
			var got []string
			if err := json.Unmarshal(stdout.Bytes(), &got); err != nil {
				t.Fatalf("decode: %v; output: %s", err, stdout)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("got %#v; want %#v", got, want)
			}
		})
	}
}

func TestWindowsBazelCanceledBeforeStart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "wrapper.cmd")
	if err := os.WriteFile(path, []byte("@echo should-not-run\r\n"), 0600); err != nil {
		t.Fatal(err)
	}
	old := *bazelPathFlag
	*bazelPathFlag = path
	defer func() { *bazelPathFlag = old }()
	b := &bazel{}
	stdout, _ := b.newCommand("version")
	b.Cancel()
	if err := b.cmd.Run(); err == nil || stdout.Len() != 0 {
		t.Fatalf("canceled command ran: %v %s", err, stdout)
	}
}

func TestWindowsBatchCancellation(t *testing.T) {
	executable, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	t.Setenv("IBAZEL_ARGV_HELPER", "1")
	t.Setenv("IBAZEL_TEST_EXE", executable)
	t.Setenv("IBAZEL_HELPER_SLEEP", "1")
	path := filepath.Join(t.TempDir(), "wrapper.cmd")
	if err := os.WriteFile(path, []byte("@\"%IBAZEL_TEST_EXE%\" %*\r\n"), 0600); err != nil {
		t.Fatal(err)
	}
	old := *bazelPathFlag
	*bazelPathFlag = path
	defer func() { *bazelPathFlag = old }()
	b := &bazel{}
	b.newCommand("-test.run=TestWindowsBazelHelper", "--")
	defer b.Cancel()
	b.cmd.Stdout = nil
	stdout, err := b.cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := b.cmd.Start(); err != nil {
		t.Fatal(err)
	}
	ready := make(chan string, 1)
	go func() { line, _ := bufio.NewReader(stdout).ReadString('\n'); ready <- line }()
	var line string
	select {
	case line = <-ready:
	case <-time.After(10 * time.Second):
		t.Fatal("helper did not start")
	}
	var pid int
	if err := json.Unmarshal([]byte(line), &pid); err != nil {
		t.Fatal(err)
	}
	process, err := syscall.OpenProcess(syscall.SYNCHRONIZE, false, uint32(pid))
	if err != nil {
		t.Fatal(err)
	}
	defer syscall.CloseHandle(process)
	b.Cancel()
	done := make(chan error, 1)
	go func() { done <- b.cmd.Wait() }()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("cancellation succeeded unexpectedly")
		}
	case <-time.After(10 * time.Second):
		t.Fatal("wrapper did not stop")
	}
	event, err := syscall.WaitForSingleObject(process, 1000)
	if err != nil || event != syscall.WAIT_OBJECT_0 {
		t.Fatalf("helper %d survived cancellation: %d, %v", pid, event, err)
	}
}

func TestWindowsBatchRejectsNewlines(t *testing.T) {
	path := filepath.Join(t.TempDir(), "wrapper.cmd")
	if err := os.WriteFile(path, []byte("@echo should-not-run\r\n"), 0600); err != nil {
		t.Fatal(err)
	}
	old := *bazelPathFlag
	*bazelPathFlag = path
	defer func() { *bazelPathFlag = old }()
	for _, arg := range []string{"line\nbreak", "line\rbreak"} {
		b := &bazel{}
		b.newCommand("query", arg)
		defer b.Cancel()
		if err := b.cmd.Run(); err == nil || b.cmd.Process != nil {
			t.Fatalf("invalid batch argument was executed: %v", err)
		}
	}
}
