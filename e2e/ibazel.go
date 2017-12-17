package e2e

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"testing"
	"time"

	bazel "github.com/bazelbuild/bazel-integration-testing/go"
)

// Maximum amount of time to wait before failing a test for not matching your expectations.
var delay = 10 * time.Second

type IBazelTester struct {
	bazel *bazel.TestingBazel
	t     *testing.T

	cmd          *exec.Cmd
	stderrBuffer *bytes.Buffer
	stdoutBuffer *bytes.Buffer
	stdoutOld    string
}

func NewIBazelTester(t *testing.T, bazel *bazel.TestingBazel) *IBazelTester {
	return &IBazelTester{
		bazel: bazel,
		t:     t,
	}
}

func (i *IBazelTester) bazelPath() string {
	return i.bazel.GetBazel()
}

func (i *IBazelTester) Run(target string) {
	i.cmd = exec.Command(ibazelPath, "--bazel_path="+i.bazelPath(), "--log_to_file=/tmp/ibazel_output.log", "run", target)

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

func (i *IBazelTester) GetOutput() string {
	return string(i.stdoutBuffer.Bytes())
}

func (i *IBazelTester) ExpectOutput(want string) {
	stopAt := time.Now().Add(delay)
	for time.Now().Before(stopAt) {
		time.Sleep(5 * time.Millisecond)

		// Grab the output and strip output that was available last time we passed
		// a test.
		out := i.GetOutput()[len(i.stdoutOld):]
		if match, err := regexp.MatchString(want, out); match == true && err == nil {
			// Save the current output value for the next iteratinog.
			i.stdoutOld = i.GetOutput()
			return
		}
	}

	if match, err := regexp.MatchString(want, i.GetOutput()); match == false || err != nil {
		i.t.Errorf("Expected iBazel output after %v to be:\nWanted [%v], got [%v]", delay, want, i.GetOutput())
		debug.PrintStack()

		// In order to prevent cascading errors where the first result failing to
		// match ruins the error output for the rest of the runs, persist the old
		// stdout.
		i.stdoutOld = i.GetOutput()
	}
}

func (i *IBazelTester) GetError() string {
	return string(i.stderrBuffer.Bytes())
}

func (i *IBazelTester) GetSubprocessPid() int64 {
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

func (i *IBazelTester) Kill() {
	if err := i.cmd.Process.Kill(); err != nil {
		panic(err)
	}
}
