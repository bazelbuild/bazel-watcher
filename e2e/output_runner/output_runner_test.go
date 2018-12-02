package output_runner

import (
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"testing"

	bazel "github.com/bazelbuild/bazel-integration-testing/go"
	"github.com/bazelbuild/bazel-watcher/e2e"
)

func must(t *testing.T, e error) {
	if e != nil {
		t.Errorf("Error: %s\n", e)
		t.Logf("Stack:\n%s", string(debug.Stack()))
	}
}

func checkNoSentinel(t *testing.T, sentinelFile *os.File, msg string) {
	if _, err := os.Stat(sentinelFile.Name()); !os.IsNotExist(err) {
		must(t, fmt.Errorf("Found a sentinel when expecting none: %s\n", msg))
	}
}

func checkSentinel(t *testing.T, sentinelFile *os.File, msg string) {
	if _, err := os.Stat(sentinelFile.Name()); os.IsNotExist(err) {
		t.Error(err)
		must(t, fmt.Errorf("Couldn't find a sentinel: %s\n%s\n", msg, err))
	}

	os.Remove(sentinelFile.Name())
}

func TestOutputRunner(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}

	sentinelFile, err := ioutil.TempFile("", "fixCommandSentinel")
	must(t, err)
	must(t, sentinelFile.Close())
	checkSentinel(t, sentinelFile, "ioutil.TempFile creates the file by default. Delete it.")
	checkNoSentinel(t, sentinelFile, "The sentinal should now be deleted.")

	must(t, b.ScratchFile(".bazel_fix_commands.json", fmt.Sprintf(`
[{
	"regex": "^(.*)runacommand(.*)$",
	"command": "touch",
	"args": ["%s"]
}]
`, sentinelFile.Name())))
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchFileWithMode("test.sh", `printf "action"`, 0777))
	must(t, b.ScratchFile("defs.bzl", `
def doit():
  print("runacommand")
`))
	must(t, b.ScratchFile("BUILD", `
load("//:defs.bzl", "doit")

doit()

sh_binary(
  name = "test",
  srcs = ["test.sh"],
)
`))

	ibazel := e2e.NewIBazelTester(t, b)
	ibazel.RunWithBazelFixCommands("//:test")

	ibazel.ExpectOutput("action")
	checkSentinel(t, sentinelFile, "The first run should create a sentinel.")

	ibazel.Kill()

	// TODO: Running the test a second time fails. I think there is a bug in the way
	// buffers are registered and they are lost between runs. Interestingly it
	// works if you reinvoke ibazel.
	ibazel = e2e.NewIBazelTester(t, b)
	ibazel.RunWithBazelFixCommands("//:test")

	ibazel.ExpectOutput("action")
	checkSentinel(t, sentinelFile, "The first run should create a sentinel.")

	ibazel.Kill()

	//// Test that the command is run again.
	//must(t, b.ScratchFileWithMode("test.sh", `printf "action"`, 0777))

	//ibazel.ExpectOutput("action2")
	//checkSentinel(t, sentinelFile, "The second run should create a sentinel.")

	// Now remove the print and it shouldn't fire.
	must(t, b.ScratchFile("defs.bzl", `
def doit():
  print("not it")
`))

	ibazel.ExpectOutput("action")
	checkNoSentinel(t, sentinelFile, "The third run should not create a sentinel.")
}

func TestNotifyWhenInvalidConfig(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}

	must(t, b.ScratchFile(".bazel_fix_commands.json", `
invalid json file
`))
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchFileWithMode("test.sh", `printf "Hello world"`, 0777))
	must(t, b.ScratchFile("BUILD", `
sh_binary(
	name = "test",
	srcs = ["test.sh"],
)
`))

	ibazel := e2e.NewIBazelTester(t, b)
	ibazel.RunWithBazelFixCommands("//:test")
	defer ibazel.Kill()

	// It should run the program and print out an error that says your JSON is
	// invalid.
	ibazel.ExpectIBazelError("Error in .bazel_fix_commands.json")
	ibazel.ExpectOutput("Hello world")
}
