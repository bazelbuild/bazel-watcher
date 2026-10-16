package new_files

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
)

const mainFiles = `
-- BUILD.bazel --
sh_binary(
  name = "contents",
  srcs = ["contents.sh"],
  data = glob(["data/**/*.txt"]),
)
-- contents.sh --
# Windows has no runfiles tree without --enable_runfiles, so read the sources.
cd "$BUILD_WORKSPACE_DIRECTORY"
printf 'contents:'
for f in $(find data -name '*.txt' | sort); do
  printf ' %s' "$(cat "$f")"
done
printf '\n'
-- data/a.txt --
a
`

func TestMain(m *testing.M) {
	e2e.TestMain(m, e2e.Args{
		Main: mainFiles,
	})
}

func TestRunWithNewFile(t *testing.T) {
	defer os.Remove(filepath.Join("data", "b.txt"))

	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//:contents")
	defer ibazel.Kill()
	ibazel.ExpectOutput("contents: a\n")

	e2e.MustWriteFile(t, filepath.Join("data", "b.txt"), "b")
	ibazel.ExpectOutput("contents: a b\n")

	e2e.MustWriteFile(t, filepath.Join("data", "b.txt"), "b2")
	ibazel.ExpectOutput("contents: a b2\n")
}

func TestRunWithNewDirectory(t *testing.T) {
	if runtime.GOOS == "darwin" {
		t.Skip("FSEvents does not report directory creation")
	}
	defer os.RemoveAll(filepath.Join("data", "sub"))

	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//:contents")
	defer ibazel.Kill()
	ibazel.ExpectOutput("contents: a\n")

	e2e.MustMkdir(t, filepath.Join("data", "sub", "deeper"))
	e2e.MustWriteFile(t, filepath.Join("data", "sub", "deeper", "c.txt"), "c")
	ibazel.ExpectOutput("contents: a c\n")

	e2e.MustWriteFile(t, filepath.Join("data", "sub", "deeper", "c.txt"), "c2")
	ibazel.ExpectOutput("contents: a c2\n")
}

func TestRunIgnoresEditorTempFiles(t *testing.T) {
	defer os.Remove(filepath.Join("data", ".b.txt.swp"))

	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//:contents")
	defer ibazel.Kill()
	ibazel.ExpectOutput("contents: a\n")

	e2e.MustWriteFile(t, filepath.Join("data", ".b.txt.swp"), "swap")
	ibazel.ExpectNoOutput(5 * time.Second)
}
