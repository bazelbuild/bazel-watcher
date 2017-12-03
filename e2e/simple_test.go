package e2e

import (
	"testing"
	"time"

	bazel "github.com/bazelbuild/bazel-integration-testing/go"
)

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

	ibazel := IBazelTester(b)
	ibazel.Run("//:test")
	defer ibazel.Kill()
	time.Sleep(2 * time.Second)
	res := ibazel.GetOutput()

	assertEqual(t, "Started!", res, "Ouput was inequal")
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

	ibazel := IBazelTester(b)
	ibazel.Run("//:test")
	defer ibazel.Kill()

	expectedOut := ""
	verify := func(startedString string) {
		expectedOut += startedString
		time.Sleep(5 * time.Second)
		assertEqual(t, expectedOut, ibazel.GetOutput(), "Ouput was inequal")
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
