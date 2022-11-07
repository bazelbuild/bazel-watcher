package simple

import (
	"fmt"
	"math"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

// Create 4096 + 1 (1 above the fsevents limit) data files in individual directories to be watched.
var dirCount = int(math.Pow(2, 12)) + 1

func TestMain(m *testing.M) {

	// Create directory structure of the form:
	//   //watched/BAZEL.build
	//   //watched/many_dirs.sh
	//   //watched/dir_[0-<dircount - 1>]/data.txt
	//	 //unwatched/data.txt

	dataFiles := make([]string, dirCount)
	dataFileNames := make([]string, dirCount)
	for i := 0; i < dirCount; i++ {
		dataFileName := fmt.Sprintf("dir_%d/data.txt", i)
		dataFileNames[i] = fmt.Sprintf("\"%s\"", dataFileName)
		if i == 0 {
			dataFiles[i] = fmt.Sprintf("-- watched/%s --\nfirst file!", dataFileName)
		} else {
			dataFiles[i] = fmt.Sprintf("-- watched/%s --", dataFileName)
		}
	}

	// Create a project that cats the contents of all data files

	mainFiles := fmt.Sprintf(`
-- watched/BUILD.bazel --

sh_binary(
  name = "many_dirs",
  srcs = ["many_dirs.sh"],
  data = [%s]
)

-- watched/many_dirs.sh --
cat watched/dir_*/data.txt

-- unwatched/data.txt --
nothing to see here

`, strings.Join(dataFileNames, ", ")) + strings.Join(dataFiles, "\n")

	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: mainFiles,
	})
}

func TestManyDirsRunWithModifiedFile(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//watched:many_dirs")
	defer ibazel.Kill()

	ibazel.ExpectOutput("first file!")

	e2e.MustWriteFile(t, "watched/dir_10/data.txt", "10th file!")
	ibazel.ExpectOutput("10th file!")

	lastFile := dirCount - 1
	e2e.MustWriteFile(t, fmt.Sprintf("watched/dir_%d/data.txt", lastFile), "last file!")
	ibazel.ExpectOutput("last file!")
}

func TestManyDirsDoesNotWatchOutsideCone(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//watched:many_dirs")
	defer ibazel.Kill()

	ibazel.ExpectOutput("first file!")

	// Give it time to start up and query.
	e2e.MustWriteFile(t, "watched/dir_10/data.txt", "10th file!")
	ibazel.ExpectOutput("10th file!")

	e2e.MustWriteFile(t, "unwatched/data.txt", "something else")
	ibazel.ExpectNoOutput(1 * time.Second)
}
