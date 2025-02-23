package simple

import (
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
)

const mainFiles = `
-- single/BUILD.bazel --
# Create an sh_test that passes and prints some output. Confirm that the
# results were streamed.
sh_test(
  name = "stream",
  srcs = ["stream.sh"],
)
-- single/stream.sh --
printf "test output"
exit 0
`

func TestMain(m *testing.M) {
	e2e.TestMain(m, e2e.Args{
		Main: mainFiles,
	})
}

func TestSimpleTest(t *testing.T) {
	// When iBazel can detect that you're testing a single target, it will run
	// the test with "--test_output=streamed" so that you can see the results
	// live.
	e2e.MustWriteFile(t, "single/stream.sh", `printf "TestSimpleTest1"`)

	ibazel := e2e.SetUp(t)
	ibazel.Test([]string{}, "//single:stream")
	defer ibazel.Kill()

	ibazel.ExpectOutput("TestSimpleTest1", 30*time.Second)

	// Now when the file is updated it should still be run in streaming mode.
	e2e.MustWriteFile(t, "single/stream.sh", `printf "TestSimpleTest2"`)
	ibazel.ExpectOutput("TestSimpleTest2", 30*time.Second)
}

func TestSingleQueryTarget(t *testing.T) {
	// When iBazel can detect that you're testing a single target, it will run
	// the test with "--test_output=streamed" so that you can see the results
	// live.

	e2e.MustWriteFile(t, "single/stream.sh", `printf "TestSingleQueryTarget"`)
	ibazel := e2e.SetUp(t)
	ibazel.Test([]string{}, "//single:all")
	defer ibazel.Kill()

	ibazel.ExpectOutput("TestSingleQueryTarget", 30*time.Second)
}

func TestMultipleQueryTarget(t *testing.T) {
	// When iBazel can detect that you're testing a single target, it will run
	// the test with "--test_output=streamed" so that you can see the results
	// live.

	e2e.MustWriteFile(t, "single/stream.sh", `printf "TestMultipleQueryTarget"`)
	ibazel := e2e.SetUp(t)
	ibazel.Test([]string{}, "//single:all", "//...")
	defer ibazel.Kill()

	ibazel.ExpectOutput("TestMultipleQueryTarget", 30*time.Second)
}

func TestExplicitlySetOutputToSummary(t *testing.T) {
	// When iBazel can detect that you're testing a single target, it will run
	// the test with "--test_output=streamed" so that you can see the results
	// live.

	e2e.MustWriteFile(t, "single/stream.sh", `printf "TestExplicitlySetOutputToSummary"`)
	ibazel := e2e.SetUp(t)
	ibazel.Test([]string{"--test_output=summary"}, "//single:all", "//...")
	defer ibazel.Kill()

	// Wait for it to pass.
	ibazel.ExpectOutput("PASSED", 30*time.Second)

	// Now confirm that the sentinel value isn't in the output.
	if strings.Contains(ibazel.GetOutput(), "TestExplicitlySetOutputToSummary") {
		t.Errorf("Wanted TestExplicitlySetOutputToSummary to not be printed, but it was. Got:\n\nOutput:\n%s", ibazel.GetOutput())
	}
}
