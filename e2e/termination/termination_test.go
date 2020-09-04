package termination

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"testing"

	"github.com/bazelbuild/bazel-watcher/e2e"
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
		SetUp: func() error {
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			if err := filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if strings.HasSuffix(path, ".sh") {
					if err := os.Chmod(path, 0777); err != nil {
						return fmt.Errorf("Error os.Chmod(%q, 0777): %v", path, err)
					}
				}
				return nil
			}); err != nil {
				fmt.Printf("Error walking dir: %v\n", err)
				return err
			}
			return nil
		},
	})
}

func TestTerminationBasic(t *testing.T) {
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
	ibazel := e2e.SetUp(t)
	e2e.MustWriteFile(t, "termination.sh", "printf \"Started 1!\";"+signalHandlerBroken)
	ibazel.Run([]string{}, "//:termination")
	ibazel.ExpectOutput("Started 1!")
	e2e.MustWriteFile(t, "termination.sh", "printf \"Started 2!\";"+signalHandlerBroken)
	ibazel.Signal(syscall.SIGINT)
	ibazel.ExpectOutput("Started 1!Started 2!")
	defer ibazel.Kill()
}
