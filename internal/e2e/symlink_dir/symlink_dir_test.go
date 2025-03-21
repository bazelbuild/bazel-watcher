package symlink_dir

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

// bazel_testing.TestMain automatically creates a `WORKSPACE` file at the root if not provided
const mainFiles = `
-- BUILD.bazel --
sh_binary(
  name = "simple",
  srcs = ["simple.sh"],
)

-- simple.sh --
printf "Started 1!"
`

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: mainFiles,
		SetUp: func() error {
			// creates a directory `./holder` and symlink
			// `./holder/symlinked-workspace` pointing to the normal working
			// directory `./main` created by `bazel_testing.TestMain`

			// create holder directory
			holder, err := filepath.Abs(filepath.Join("..", "holder"))
			if err != nil {
				log.Fatalf("Error determining absolute path: %v", err)
			}

			if err = os.Mkdir(holder, 0777); err != nil {
				log.Fatalf("Error making directory: %v", err)
			}

			// create symlink pointing to the main workspace
			symlinkWd, err := filepath.Abs(filepath.Join(holder, "symlinked-workspace"))
			if err != nil {
				log.Fatalf("Error determining absolute path: %v", err)
			}

			cwd, err := os.Getwd()
			if err != nil {
				log.Fatalf("Error getting working directory: %v", err)
			}

			if err = os.Symlink(cwd, symlinkWd); err != nil {
				log.Fatalf("Error creating symlink: %v", err)
			}

			// chdir via the symlink
			if err = os.Chdir(symlinkWd); err != nil {
				log.Fatalf("Error changing directory: %v", err)
			}

			return err
		},
	})
}

// note: the implementation of `workspace.go` uses `os.Getwd` which states "If
// the current directory can be reached via multiple paths (due to symbolic
// links), Getwd may return any one of them." This resulted in inconsistent
// automatic rebuild behavior because `ibazel.go` was not evaluating symlinks.
func TestSymlinkRun(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//:simple")
	defer ibazel.Kill()

	ibazel.ExpectOutput("Started 1!", 35 * time.Second)

	e2e.MustWriteFile(t, "simple.sh", `printf "Started 2!"`)
	ibazel.ExpectOutput("Started 2!")
}
