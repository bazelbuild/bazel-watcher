package output_runner

import (
	"os"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
)

const mainFiles = `
-- single/defs.bzl --
def fix_deps():
  print("runacommand")
-- single/BUILD --
load("//single:defs.bzl", "fix_deps")

fix_deps()

sh_binary(
  name = "test",
  srcs = ["test.sh"],
)

sh_binary(
  name = "overwrite",
  srcs = ["overwrite.sh"],
)
-- single/test.sh --
printf "action"
-- single/overwrite.sh --
printf "overwrite"
-- multiple/defs.bzl --
def fix_deps():
  print("runcommand foo")
  print("runcommand bar")
  print("runcommand foo")
  print("runcommand baz")
-- multiple/BUILD --
load("//multiple:defs.bzl", "fix_deps")

fix_deps()

sh_binary(
  name = "test",
  srcs = ["test.sh"],
)
-- multiple/test.sh --
printf "action"
`

func TestMain(m *testing.M) {
	e2e.TestMain(m, e2e.Args{
		Main: mainFiles,
	})
}

func checkNoSentinel(t *testing.T, sentinelFile *os.File, msg string) {
	t.Helper()

	if _, err := os.Stat(sentinelFile.Name()); !os.IsNotExist(err) {
		t.Errorf("Found a sentinel when expecting none: %s\n", msg)
	}
}

func checkSentinel(t *testing.T, sentinelFile *os.File, msg string) {
	t.Helper()
	sentinalFileName := sentinelFile.Name()

	deadline := time.Now().Add(5 * time.Second)
	var err error
	for {
		if time.Now().After(deadline) {
			t.Errorf("Couldn't find sentinal. os.Stat(%q): %s\n%s\n", sentinalFileName, err, msg)
			return
		} else if _, err := os.Stat(sentinalFileName); err == nil {
			// No error stat'ing the file means it exists.
			os.Remove(sentinelFile.Name())
			return
		}
	}
}

func TestNotifyWhenInvalidConfig(t *testing.T) {
	e2e.MustWriteFile(t, ".bazel_fix_commands.json", `
invalid json file
`)

	ibazel := e2e.SetUp(t)
	ibazel.RunWithBazelFixCommands("//single:test")
	defer ibazel.Kill()

	// It should run the program and print out an error that says your JSON is
	// invalid.
	ibazel.ExpectIBazelError("Error in .bazel_fix_commands.json")
	ibazel.ExpectOutput("action")
}
