package e2e

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	bazel "github.com/bazelbuild/bazel-integration-testing/go"
)

type iBazelTester struct {
	bazel *bazel.TestingBazel

	cmd          *exec.Cmd
	stderrBuffer *bytes.Buffer
	stdoutBuffer *bytes.Buffer
}

func IBazelTester(bazel *bazel.TestingBazel) *iBazelTester {
	return &iBazelTester{
		bazel: bazel,
	}
}

func (i *iBazelTester) bazelPath() string {
	return i.bazel.GetBazel()
}

func (i *iBazelTester) Run(target string) {
	i.cmd = exec.Command(ibazelPath, "--bazel_path="+i.bazelPath(), "--log_to_file=/tmp/output.log", "run", target)

	errCode, buildStdout, buildStderr := i.bazel.RunBazel([]string{"build", target})
	if errCode != 0 {
		panic(fmt.Sprintf("Unable to build target. Error code: %d\nStdout:\n%s\nStderr:\n%s", errCode, buildStdout, buildStderr))
	}

	i.stdoutBuffer = &bytes.Buffer{}
	i.cmd.Stdout = i.stdoutBuffer

	i.stderrBuffer = &bytes.Buffer{}
	i.cmd.Stderr = i.stderrBuffer

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
