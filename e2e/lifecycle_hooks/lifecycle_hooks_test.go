package lifecycle_hooks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-watcher/e2e"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

const mainFiles = `
-- BUILD.bazel --
sh_binary(
  name = "test",
  srcs = ["test.sh"],
)
-- test.sh --
printf "action"
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

func TestLifecycleHooks(t *testing.T) {
	ibazel := e2e.SetUp(t)

	ibazel.RunWithAdditionalArgs("//:test", []string{
		"-run_command_before=echo hi-before",
		"-run_command_after=echo hi-after",
	})
	ibazel.ExpectOutput("hi-before")
	ibazel.ExpectOutput("hi-after")
	ibazel.Kill()
}
