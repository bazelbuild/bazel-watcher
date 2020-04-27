package error

import (
	"testing"

	"github.com/bazelbuild/bazel-watcher/e2e"
	"github.com/bazelbuild/bazel-watcher/e2e/example_client"
)

func TestMain(m *testing.M) {
	example_client.TestMain(m)
}

func TestSimpleBuildWithoutSourceFiles(t *testing.T) {
	e2e.MustWriteFile(t, "BUILD", `
# Invalid rule due to a missing input file
sh_binary(
  name = "test",
  srcs = ["test.sh"],
)
`)

	ibazel := e2e.SetUp(t)
	ibazel.Build("//:test")
	defer ibazel.Kill()

	ibazel.ExpectError("//:test: missing input file '//:test.sh'")
}

func TestSimpleBuildWithQueryFailure(t *testing.T) {
	e2e.MustWriteFile(t, "BUILD", `
# Invalid rule due to a typo
shh_binary(
  name = "test",
  srcs = ["test.sh"],
)
`)

	ibazel := e2e.SetUp(t)
	ibazel.Build("//:test")
	defer ibazel.Kill()

	ibazel.ExpectError("name 'shh_binary' is not defined")
}

/*
func TestExampleClientWhoDies(t *testing.T) {
	c := example_client.StartLiveReload(t)
	defer c.Cleanup()

	if v := c.GetRaw(t); v != "1" {
		t.Errorf("Expected raw value to be 1, got %q", v)
	}

	// Simulate a crash in the server by killing it.
	c.Kill(t)

	// Give the server a bit of time to "crash".
	time.Sleep(2 * time.Second)

	// Now restart the crashed server by changing a file that it is watching.
	// This should cause ibazel to relaunch the program that is running.
	c.SetData(t, "2")

	// Since we are restarting the server it might need a bit of time to restart.
	time.Sleep(2 * time.Second)

	// Redetect the paths that are associated with the freshly restarted
	// instance.
	c.DetectServerParameters(t)

	if v := c.GetRaw(t); v != "2" {
		t.Errorf("Expected raw value to be 2, got %q", v)
	}
}
*/
