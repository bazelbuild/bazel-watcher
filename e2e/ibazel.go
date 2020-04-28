package e2e

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

// Maximum amount of time to wait before failing a test for not matching your expectations.
const (
	defaultDelay = 20 * time.Second
)

type IBazelTester struct {
	t             *testing.T
	ibazelLogFile string

	cmd          *exec.Cmd
	stderrBuffer *Buffer
	stderrOld    string
	stdoutBuffer *Buffer
	stdoutOld    string
	ibazelErrOld string
}

func NewIBazelTester(t *testing.T) *IBazelTester {
	f, err := ioutil.TempFile("", "ibazel_output.*.log")
	if err != nil {
		panic(fmt.Sprintf("Error ioutil.Tempfile: %v", err))
	}

	return &IBazelTester{
		t:             t,
		ibazelLogFile: f.Name(),
	}
}

func (i *IBazelTester) bazelPath() string {
	i.t.Helper()
	path, err := exec.LookPath("bazel")
	if err != nil {
		i.t.Fatalf("Unable to find bazel binary: %v", err)
	}
	return path
}

func (i *IBazelTester) Build(target string) {
	i.t.Helper()
	i.build(target, []string{})
}

func (i *IBazelTester) Run(bazelArgs []string, target string) {
	i.t.Helper()
	i.run(target, bazelArgs, []string{
		"--log_to_file=" + i.ibazelLogFile,
	})
}

func (i *IBazelTester) RunWithProfiler(target string, profiler string) {
	i.t.Helper()
	i.run(target, []string{}, []string{
		"--log_to_file=" + i.ibazelLogFile,
		"--profile_dev=" + profiler,
	})
}

func (i *IBazelTester) RunWithBazelFixCommands(target string) {
	i.t.Helper()
	i.run(target, []string{}, []string{
		"--log_to_file=" + i.ibazelLogFile,
		"--run_output=true",
		"--run_output_interactive=false",
	})
}

func (i *IBazelTester) GetOutput() string {
	i.t.Helper()
	return i.stdoutBuffer.String()
}

func (i *IBazelTester) ExpectOutput(want string, delay ...time.Duration) {
	i.t.Helper()

	i.checkExit()

	d := defaultDelay
	if len(delay) == 1 {
		d = delay[0]
	}
	i.Expect(want, i.GetOutput, &i.stdoutOld, d)
}

func (i *IBazelTester) ExpectError(want string, delay ...time.Duration) {
	i.t.Helper()

	i.checkExit()

	d := defaultDelay
	if len(delay) == 1 {
		d = delay[0]
	}
	i.Expect(want, i.GetError, &i.stderrOld, d)
}

func (i *IBazelTester) ExpectIBazelError(want string, delay ...time.Duration) {
	i.t.Helper()

	i.checkExit()

	d := defaultDelay
	if len(delay) == 1 {
		d = delay[0]
	}
	i.Expect(want, i.GetIBazelError, &i.ibazelErrOld, d)
}

func (i *IBazelTester) GetIBazelError() string {
	i.t.Helper()

	i.checkExit()

	iBazelError, err := os.Open(i.ibazelLogFile)
	if err != nil {
		i.t.Errorf("Error os.Open(%q): %v", i.ibazelLogFile, err)
		return ""
	}

	b, err := ioutil.ReadAll(iBazelError)
	if err != nil {
		i.t.Fatalf("Error ioutil.ReadAll(iBazelError): %v", err)
	}

	return string(b)
}

func (i *IBazelTester) Expect(want string, stream func() string, history *string, delay time.Duration) {
	i.t.Helper()

	stopAt := time.Now().Add(delay)
	for time.Now().Before(stopAt) {
		time.Sleep(5 * time.Millisecond)

		// Grab the output and strip output that was available last time we passed
		// a test.
		out := stream()[len(*history):]
		if match, err := regexp.MatchString(want, out); match == true && err == nil {
			// Save the current output value for the next iteration.
			*history = stream()
			return
		}
	}

	if match, err := regexp.MatchString(want, stream()); match == false || err != nil {
		i.t.Errorf("Expected iBazel output after %v to be:\nWanted [%v], got [%v]", delay, want, stream())
		i.t.Errorf("Stderr: [%v]\niBazelStderr: [%v]", i.GetError(), i.GetIBazelError())
		//i.t.Log(string(debug.Stack()))

		// In order to prevent cascading errors where the first result failing to
		// match ruins the error output for the rest of the runs, persist the old
		// stdout.
		*history = stream()
	}
}

func (i *IBazelTester) GetError() string {
	i.t.Helper()
	return i.stderrBuffer.String()
}

func (i *IBazelTester) GetSubprocessPid() int64 {
	i.t.Helper()
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
	i.t.Helper()
	if err := i.cmd.Process.Kill(); err != nil {
		panic(err)
	}
}

func (i *IBazelTester) build(target string, additionalArgs []string) {
	i.t.Helper()
	args := []string{"--bazel_path=" + i.bazelPath()}
	args = append(args, additionalArgs...)
	args = append(args, "build")
	args = append(args, target)
	i.cmd = exec.Command(ibazelPath, args...)

	i.stdoutBuffer = &Buffer{}
	i.cmd.Stdout = i.stdoutBuffer

	i.stderrBuffer = &Buffer{}
	i.cmd.Stderr = i.stderrBuffer

	if err := i.cmd.Start(); err != nil {
		i.t.Fatalf("Command: %s\nError: %v", i.cmd, err)
	}
}

func (i *IBazelTester) checkExit() {
	if i.cmd != nil && i.cmd.ProcessState != nil && i.cmd.ProcessState.Exited() == true {
		i.t.Errorf("ibazel is exited")
	}
}

func (i *IBazelTester) run(target string, bazelArgs []string, additionalArgs []string) {
	i.t.Helper()

	args := []string{"--bazel_path=" + i.bazelPath()}
	args = append(args, additionalArgs...)
	args = append(args, "run")
	args = append(args, target)
	args = append(args, bazelArgs...)
	i.cmd = exec.Command(ibazelPath, args...)
	i.t.Logf("ibazel invoked as: %s", strings.Join(i.cmd.Args, " "))

	cmd := bazel_testing.BazelCmd("build", target)

	var buildStdout, buildStderr bytes.Buffer
	cmd.Stdout = &buildStdout
	cmd.Stderr = &buildStderr

	// Before doing anything crazy, let's build the target to make sure it works.

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			status := exitErr.Sys().(syscall.WaitStatus)
			i.t.Fatalf("Unable to build target. Error code: %d\nStdout:\n%s\nStderr:\n%s", status.ExitStatus(), buildStdout.String(), buildStderr.String())
		}
	}

	i.stdoutBuffer = &Buffer{}
	i.cmd.Stdout = i.stdoutBuffer

	i.stderrBuffer = &Buffer{}
	i.cmd.Stderr = i.stderrBuffer

	if err := i.cmd.Start(); err != nil {
		i.t.Fatalf("Command: %s", i.cmd)
	}
}
