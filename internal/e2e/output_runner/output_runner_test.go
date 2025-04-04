package iutput_runner

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
)

const mainFiles = `
-- WORKSPACE --
#required
-- single/defs.bzl --
def fix_deps(name):
    print("used-to-be-the-magic-command")
-- single/BUILD --
load("//single:defs.bzl", "fix_deps")

genrule(name="gen_overwrite", srcs=["overwrite.sh"], outs=["outrun1"],cmd="echo runacommand >&2 && cat $< > $@")

sh_binary(
  name = "test",
  srcs = ["test.sh"],
)

sh_binary(
  name = "overwrite",
  srcs = [":gen_overwrite"],
)
-- single/test.sh --
printf "action"
-- single/overwrite.sh --
printf "overwrite"
-- multiple/BUILD --
# maybe I can dict-comprehension later
genrule(name="fix_deps", srcs=["test.sh"], outs=["out-fixdeps"],cmd="echo runcommand foo >&2 && echo runcommand bar >&2 && echo runcommand foo >&2 && echo runcommand baz >&2 && cat $< > $@")

sh_binary(
  name = "test",
  srcs = ["test.sh"],
)
genrule(
  name = "collector",
  cmd = "echo 'date' > $@",
  executable = True,
  outs = ["out-collector"],
  srcs = [ ":fix_deps" ]
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

func TestOutputRunner(t *testing.T) {
	sentinelFile, err := ioutil.TempFile("", "fixCommandSentinel")
	if err != nil {
		t.Errorf("ioutil.TempFile(\"\", \"fixCommandSentinel\": %v", err)
	}
	sentinalFileName := strings.Replace(sentinelFile.Name(), "\\", "/", -1)

	e2e.Must(t, sentinelFile.Close())
	checkSentinel(t, sentinelFile, "ioutil.TempFile creates the file by default. Delete it.")
	checkNoSentinel(t, sentinelFile, "The sentinal should now be deleted.")

	// First check that it doesn't run if there isn't a `.bazel_fix_commands.json` file.
	ibazel := e2e.SetUp(t)
        t.Log("Running //single:overwrite with no bazel_fix_commands")
	ibazel.RunWithBazelFixCommands("//single:overwrite")

	// Ensure it prints out the banner.
	ibazel.ExpectIBazelError("Did you know", 50 * time.Second)

	ibazel.Kill()

	ibazel = e2e.SetUp(t)
        t.Log("Creating bazel_fix_commands")
	e2e.MustWriteFile(t, ".bazel_fix_commands.json", fmt.Sprintf(`
	[{
		"regex": "^(.*)runacommand(.*)$",
		"command": "touch",
		"args": ["%s"]
	}]`, sentinalFileName))

        t.Log("overwriting overwrite.sh")
	e2e.MustWriteFile(t, "single/overwrite.sh", `
printf "overwrite1"
`)

        t.Log("Running //single:overwrite with bazel_fix_commands")
	ibazel.RunWithBazelFixCommands("//single:overwrite")

	ibazel.ExpectOutput("overwrite1")
	checkSentinel(t, sentinelFile, "The first run should create a sentinel.")

	ibazel.Kill()

	// Invoke the test a 2nd time to ensure it works over multiple separate
	// invocations of ibazel.
	ibazel = e2e.SetUp(t)

	// additional change in dependency of the genrule() so that the side-effect printed text is triggered
	e2e.MustWriteFile(t, "single/overwrite.sh", `
printf "overwrite3"
`)
	ibazel.RunWithBazelFixCommands("//single:overwrite")
	ibazel.ExpectOutput("overwrite3")
	checkSentinel(t, sentinelFile, "The second run should create a sentinel.")

	// TODO: Figure out why the 2nd invocation doesn't touch the file.
	// Test that the command is run again.
	//e2e.MustWriteFile(t, "overwrite.sh", `printf "overwrite2"`)

	//ibazel.ExpectOutput("overwrite2")
	//checkSentinel(t, sentinelFile, "The third run should create a sentinel.")

	// Now replace the print and it shouldn't fire.
        // More recent trick/fact: not updating overwrite.sh means no source deps changed
	e2e.MustWriteFile(t, "defs.bzl", `
def fix_deps():
  print("not it")
`)

	ibazel.ExpectOutput("overwrite3")
	checkNoSentinel(t, sentinelFile, "The third run should not create a sentinel.")
}

func TestOutputRunnerUniqueCommandsOnly(t *testing.T) {
	e2e.MustWriteFile(t, ".bazel_fix_commands.json", `
       [{
               "regex": "^.*runcommand (.*)$",
               "command": "echo",
               "args": ["$1"]
       }]`)

	ibazel := e2e.NewIBazelTester(t)
	ibazel.RunWithBazelFixCommands("//multiple:collector")
	defer ibazel.Kill()

	ibazel.ExpectFixCommands([]string{
		"echo foo",
		"echo bar",
		"echo baz",
	})
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
