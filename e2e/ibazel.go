package e2e

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

type iBazelTester struct {
	target  string
	binPath string

	cmd          *exec.Cmd
	stderrBuffer *bytes.Buffer
	stdoutBuffer *bytes.Buffer
}

func IBazelTester(target, binPath string) *iBazelTester {
	return &iBazelTester{
		target:  target,
		binPath: binPath,
	}
}

func (i *iBazelTester) Run() {
	i.cmd = exec.Command(ibazelPath, "--bazel_path="+bazelPath, "--log_to_file=/tmp/output.log", "run", i.target)

	i.stdoutBuffer = &bytes.Buffer{}
	i.cmd.Stdout = i.stdoutBuffer

	i.stderrBuffer = &bytes.Buffer{}
	i.cmd.Stderr = i.stderrBuffer

	launcherPath := filepath.Join(os.TempDir(), "ibazel_e2e_subprocess_launcher")
	// Try to delete the file. Sometimes the file won't be overwritten properly
	// but if it is deleted there is no risk for that problem. Investigate that.
	os.Remove(launcherPath)
	launcher := fmt.Sprintf("#! /usr/bin/env bash\nexec %s", i.binPath)
	err := ioutil.WriteFile(launcherPath, []byte(launcher), 0755)
	if err != nil {
		panic(err)
	}

	if err := i.cmd.Start(); err != nil {
		fmt.Printf("Command: %s", i.cmd)
		panic(err)
	}
}

func (i *iBazelTester) GetOutput() string {
	return string(i.stdoutBuffer.Bytes())
}
func (i *iBazelTester) GetError() string {
	return string(i.stderrBuffer.Bytes())
}
func (i *iBazelTester) GetSubprocessPid() int64 {
	f, err := os.Open(filepath.Join(os.TempDir(), "ibazel_e2e_subprocess_launcher.pid"))
	if err != nil {
		panic(err)
	}

	rawPid, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	pid, err := strconv.ParseInt(string(rawPid), 10, 32)
	if err != nil {
		panic(err)
	}
	return pid
}
func (i *iBazelTester) Kill() {
	if err := i.cmd.Process.Kill(); err != nil {
		panic(err)
	}
}
