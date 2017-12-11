package simple

import (
	"os"
	"reflect"
	"runtime/debug"
	"testing"
	"time"

	bazel "github.com/bazelbuild/bazel-integration-testing/go"
	"github.com/bazelbuild/bazel-watcher/e2e"
)

func must(t *testing.T, e error) {
	if e != nil {
		t.Errorf("Error: %s", e)
		debug.PrintStack()
	}
}

func assertNotEqual(t *testing.T, want, got interface{}, msg string) {
	if reflect.DeepEqual(want, got) {
		t.Errorf("Wanted %s, got %s. %s", want, got, msg)
		debug.PrintStack()
	}
}
func assertEqual(t *testing.T, want, got interface{}, msg string) {
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Wanted [%v], got [%v]. %s", want, got, msg)
		debug.PrintStack()
	}
}

func TestSimpleRun(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchFileWithMode("test.sh", `printf "Started!"`, 0777))
	must(t, b.ScratchFile("BUILD", `
sh_binary(
	name = "test",
	srcs = ["test.sh"],
)
`))

	ibazel := e2e.NewIBazelTester(b)
	ibazel.Run("//:test")
	defer ibazel.Kill()
	time.Sleep(2 * time.Second)
	res := ibazel.GetOutput()

	assertEqual(t, "Started!", res, "Output was unequal")
}

func TestSimpleRunUnderSubdir(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchDir("subdir"))
	must(t, b.ScratchFileWithMode("subdir/test.sh", `printf "Started!"`, 0777))
	must(t, b.ScratchFile("subdir/BUILD", `
sh_binary(
	name = "test",
	srcs = ["test.sh"],
)
`))

	ibazel := e2e.NewIBazelTester(b)

	err = os.Chdir("subdir")
	if err != nil {
		t.Fatal(err)
	}

	ibazel.Run("test")
	defer ibazel.Kill()
	
	err = os.Chdir("..")
	if err != nil {
		t.Fatal(err)
	}

	time.Sleep(2 * time.Second)
	res := ibazel.GetOutput()

	assertEqual(t, "Started!", res, "Output was unequal")
}

func TestSimpleRunWithModifiedFile(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchFileWithMode("test.sh", `printf "Started!"`, 0777))
	must(t, b.ScratchFile("BUILD", `
sh_binary(
	name = "test",
	srcs = ["test.sh"],
)
`))

	ibazel := e2e.NewIBazelTester(b)
	ibazel.Run("//:test")
	defer ibazel.Kill()

	expectedOut := ""
	verify := func(startedString string) {
		expectedOut += startedString
		time.Sleep(5 * time.Second)
		assertEqual(t, expectedOut, ibazel.GetOutput(), "Output was unequal")
	}

	// Give it time to start up and query.
	must(t, b.ScratchFileWithMode("test.sh", `printf "Started2!"`, 0777))
	verify("Started2!")

	// Manipulate a source file and sleep past the debounce.
	must(t, b.ScratchFileWithMode("test.sh", `printf "Started3!"`, 0777))
	verify("Started3!")

	// Now a BUILD file.
	must(t, b.ScratchFile("BUILD", `
sh_binary(
	# New comment
	name = "test",
	srcs = ["test.sh"],
)
`))
	verify("Started3!")
}
