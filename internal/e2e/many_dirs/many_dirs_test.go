package simple

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

var dirCount = int(math.Pow(2, 12)) + 1

func TestMain(m *testing.M) {

	// Create 4096 + 1 (1 above the fsevents limit) data files in individual directories to be watched

	dataFiles := make([]string, dirCount)
	dataFileNames := make([]string, dirCount)
	for i := 0; i < dirCount; i++ {
		dataFileName := fmt.Sprintf("dir_%d/data.txt", i)
		dataFileNames[i] = fmt.Sprintf("\"%s\"", dataFileName)
		if i == 0 {
			dataFiles[i] = fmt.Sprintf("-- %s --\nfirst file!", dataFileName)
		} else {
			dataFiles[i] = fmt.Sprintf("-- %s --", dataFileName)
		}
	}

	// Create a project that cats the contents of all data files

	mainFiles := fmt.Sprintf(`
-- BUILD.bazel --

sh_binary(
  name = "many_dirs",
  srcs = ["many_dirs.sh"],
  data = [%s]
)

-- many_dirs.sh --
cat dir_*/data.txt

`, strings.Join(dataFileNames, ", ")) + strings.Join(dataFiles, "\n")

	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: mainFiles,
	})
}

func TestManyDirsRunWithModifiedFile(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//:many_dirs")
	defer ibazel.Kill()

	ibazel.ExpectOutput("first file!")

	e2e.MustWriteFile(t, "dir_10/data.txt", "10th file!")
	ibazel.ExpectOutput("10th file!")

	lastFile := dirCount - 1
	e2e.MustWriteFile(t, fmt.Sprintf("dir_%d/data.txt", lastFile), "last file!")
	ibazel.ExpectOutput("last file!")
}
