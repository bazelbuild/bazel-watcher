package e2e

import (
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
var delay = 20 * time.Second

type IBazelTester struct {
	bazel *bazel.TestingBazel
	t     *testing.T

	cmd          *exec.Cmd
	stderrBuffer *Buffer
	stderrOld    string
	stdoutBuffer *Buffer
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

func (i *IBazelTester) Build(target string) {
	i.build(target, []string{})
}

func (i *IBazelTester) Run(target string) {
	i.run(target, []string{})
}

func (i *IBazelTester) RunWithProfiler(target string, profiler string) {
	i.run(target, []string{"--profile_dev=" + profiler})
}

func (i *IBazelTester) RunWithBazelFixCommands(target string) {
	i.run(target, []string{
		"--run_output=true",
		"--run_output_interactive=false",
	})
}

func (i *IBazelTester) GetOutput() string {
	return i.stdoutBuffer.String()
}

func (i *IBazelTester) ExpectOutput(want string) {
	i.Expect(want, i.GetOutput, &i.stdoutOld)
}

func (i *IBazelTester) ExpectError(want string) {
	i.Expect(want, i.GetError, &i.stderrOld)
}

func (i *IBazelTester) ExpectIBazelError(want string) {
}

func (i *IBazelTester) GetIBazelError() string {
	iBazelError, err := os.Open("/tmp/ibazel_output.log")

	b, err := ioutil.ReadAll(iBazelError)
	if err != nil {
		i.t.Fatal(err)
	}

	return string(b)
}

func (i *IBazelTester) Expect(want string, stream func() string, history *string) {
	stopAt := time.Now().Add(delay)
	for time.Now().Before(stopAt) {
		time.Sleep(500 * time.Millisecond)

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
		i.t.Log(string(debug.Stack()))

		// In order to prevent cascading errors where the first result failing to
		// match ruins the error output for the rest of the runs, persist the old
		// stdout.
		*history = stream()
	}
}

func (i *IBazelTester) GetError() string {
	return i.stderrBuffer.String()
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

func (i *IBazelTester) build(target string, additionalArgs []string) {
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
		i.t.Logf("Command: %s", i.cmd)
		panic(err)
	}
}

func (i *IBazelTester) run(target string, additionalArgs []string) {
	args := []string{"--bazel_path=" + i.bazelPath()}
	args = append(args, additionalArgs...)
	args = append(args, "run")
	args = append(args, target)
	i.cmd = exec.Command(ibazelPath, args...)

	errCode, buildStdout, buildStderr := i.bazel.RunBazel([]string{"build", target})
	if errCode != 0 {
		i.t.Fatalf("Unable to build target. Error code: %d\nStdout:\n%s\nStderr:\n%s", errCode, buildStdout, buildStderr)
	}

	i.stdoutBuffer = &Buffer{}
	i.cmd.Stdout = i.stdoutBuffer

	i.stderrBuffer = &Buffer{}
	i.cmd.Stderr = i.stderrBuffer

	if err := i.cmd.Start(); err != nil {
		i.t.Fatalf("Command: %s", i.cmd)
	}
}
