package simple

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"

	"github.com/bazelbuild/bazel-watcher/e2e"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

const mainFiles = `
-- BUILD.bazel --
# Base case test
sh_binary(
  name = "simple",
  srcs = ["simple.sh"],
)

# Environment variable tests
sh_binary(
  name = "environment",
  srcs = ["environment.sh"],
)

# --define tests
config_setting(
  name = "test_is_2",
  values = {"define": "test_number=2"},
)

sh_binary(
  name = "define",
  srcs = select({
        ":test_is_2": ["define_test_2.sh"],
        "//conditions:default": ["define_test_1.sh"],
    }),
)
-- simple.sh --
printf "Started!"
-- environment.sh --
printf "Started and IBAZEL=${IBAZEL}!"
-- define_test_1.sh --
printf "define_test_1"
-- define_test_2.sh --
printf "define_test_2"
-- subdir/BUILD.bazel --
sh_binary(
  name = "subdir",
  srcs = ["subdir.sh"],
)
-- subdir/subdir.sh --
printf "Hello subdir!"
`

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

func TestSimpleBuild(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//:simple")
	defer ibazel.Kill()

	ibazel.ExpectOutput("Started!")
}

func TestSimpleRunAfterShutdown(t *testing.T) {
	cmd := bazel_testing.BazelCmd("shutdown")
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			status := exitErr.Sys().(syscall.WaitStatus)
			if status.ExitStatus() != 0 {
				t.Fatal(errors.New("bazel failed to shut down"))
			}
		}
	}

	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//:simple")
	defer ibazel.Kill()

	ibazel.ExpectOutput("Started!")
}

func TestSimpleRunConfirmEnvironment(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//:environment")
	defer ibazel.Kill()

	ibazel.ExpectOutput("Started and IBAZEL=true!")
}

func TestSimpleRunUnderSubdir(t *testing.T) {
	subdir := "subdir"

	ibazel := e2e.SetUp(t)

	err := os.Chdir(subdir)
	if err != nil {
		t.Fatalf("Error os.Chdir(%q): %v", subdir, err)
	}
	defer func() {
		err := os.Chdir("..")
		if err != nil {
			t.Fatalf("Error os.Chdir(\"..\"): %v", err)
		}
	}()

	ibazel.Run([]string{}, "//subdir")
	defer ibazel.Kill()

	ibazel.ExpectOutput("Hello subdir")
}

func TestSimpleRunWithModifiedFile(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//:simple")
	defer ibazel.Kill()

	ibazel.ExpectOutput("Started!")

	// Give it time to start up and query.
	e2e.MustWriteFile(t, "simple.sh", `printf "Started2!"`)
	ibazel.ExpectOutput("Started2!")

	// Manipulate a source file and sleep past the debounce.
	e2e.MustWriteFile(t, "simple.sh", `printf "Started3!"`)
	ibazel.ExpectOutput("Started3!")

	// TODO: put these in directories instead of storing the old value and rewriting it
	oldValue, err := ioutil.ReadFile("BUILD.bazel")
	if err != nil {
		t.Errorf("Unable to Readfile(\"BUILD.bazel\"): %v", err)
	}
	defer e2e.MustWriteFile(t, "BUILD.bazel", string(oldValue))
	// END TODO

	// Now a BUILD.bazel file.
	e2e.MustWriteFile(t, "BUILD.bazel", `
sh_binary(
	# New comment
	name = "test",
	srcs = ["test.sh"],
)
`)
	ibazel.ExpectOutput("Started3!")
}

func TestSimpleRunWithFlag(t *testing.T) {
	ibazel := e2e.SetUp(t)

	ibazel.Run([]string{}, "//:define")
	ibazel.ExpectOutput("define_test_1")
	ibazel.Kill()

	ibazel = e2e.NewIBazelTester(t)
	ibazel.Run([]string{"--define=test_number=2"}, "//:define")
	ibazel.ExpectOutput("define_test_2")
	ibazel.Kill()

	ibazel = e2e.NewIBazelTester(t)
	ibazel.Run([]string{}, "//:define")
	ibazel.ExpectOutput("define_test_1")
	ibazel.Kill()
}
