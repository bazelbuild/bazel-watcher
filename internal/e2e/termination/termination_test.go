package termination

import (
	"runtime"
	"syscall"
	"testing"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

const signalHandler = `
tail -f /dev/null & PID=$!
_handler() { printf "$1."; kill $PID; }
trap '_handler SIGTERM' TERM
wait
`

const signalHandlerBroken = `
trap -- '' SIGTERM
tail -f /dev/null
`

const mainFiles = `
-- BUILD.bazel --
sh_binary(
  name = "termination",
  srcs = ["termination.sh"],
)
-- termination.sh --
` + signalHandler

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: mainFiles,
	})
}

func TestTerminationBasic(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("termination tests are currently broken on Windows. We would love to test this but don't have a box to test on.")
	}

	ibazel := e2e.SetUp(t)
	e2e.MustWriteFile(t, "termination.sh", "printf \"Started 1!\";"+signalHandler)
	ibazel.Run([]string{}, "//:termination")
	ibazel.ExpectOutput("Started 1!")
	e2e.MustWriteFile(t, "termination.sh", "printf \"Started 2!\";"+signalHandler)

	// Windows doesn't support signals unfortunately
	if runtime.GOOS == "windows" {
		ibazel.ExpectOutput("Started 1!Started 2!")
		ibazel.Signal(syscall.SIGINT)
		ibazel.ExpectOutput("Started 1!Started 2!")
	} else {
		ibazel.ExpectOutput("Started 1!SIGTERM.Started 2!")
		ibazel.Signal(syscall.SIGINT)
		ibazel.ExpectOutput("Started 1!SIGTERM.Started 2!SIGTERM.")
	}
	defer ibazel.Kill()
}

func TestTerminationTimeout(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("termination tests are currently broken on Windows. We would love to test this but don't have a box to test on.")
	}

	ibazel := e2e.SetUp(t)
	e2e.MustWriteFile(t, "termination.sh", "printf \"Started 1!\";"+signalHandlerBroken)
	ibazel.Run([]string{}, "//:termination")
	ibazel.ExpectOutput("Started 1!")
	e2e.MustWriteFile(t, "termination.sh", "printf \"Started 2!\";"+signalHandlerBroken)
	ibazel.Signal(syscall.SIGINT)
	ibazel.ExpectOutput("Started 1!Started 2!")
	defer ibazel.Kill()
}
