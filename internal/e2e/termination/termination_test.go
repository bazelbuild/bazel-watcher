package termination

import (
	"runtime"
	"syscall"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
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
	e2e.TestMain(m, e2e.Args{
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
	ibazel.ExpectOutput("Started 1!", 50 * time.Second)
	e2e.MustWriteFile(t, "termination.sh", "printf \"Started 2!\";"+signalHandler)

	// Windows doesn't support signals unfortunately
	if runtime.GOOS == "windows" {
		ibazel.ExpectOutput("Started 1!Started 2!", 50 * time.Second)
		ibazel.Signal(syscall.SIGINT)
		ibazel.ExpectOutput("Started 1!Started 2!", 50 * time.Second)
	} else {
		ibazel.ExpectOutput("Started 1!SIGTERM.Started 2!", 50 * time.Second)
		ibazel.Signal(syscall.SIGINT)
		ibazel.ExpectOutput("Started 1!SIGTERM.Started 2!SIGTERM.", 50 * time.Second)
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
	ibazel.ExpectOutput("Started 1!", 50 * time.Second)
	e2e.MustWriteFile(t, "termination.sh", "printf \"Started 2!\";"+signalHandlerBroken)
	ibazel.Signal(syscall.SIGINT)
	ibazel.ExpectOutput("Started 1!Started 2!", 50 * time.Second)
	defer ibazel.Kill()
}
