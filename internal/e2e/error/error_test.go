package error

import (
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
	"github.com/bazelbuild/bazel-watcher/internal/e2e/example_client"
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

	ibazel.ExpectError("//:test: missing input file '//:test.sh'", 35 * time.Second)
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

func TestExampleClientWhoDies(t *testing.T) {
	e2e.MustWriteFile(t, "BUILD", `
sh_binary(
	name = "live_reload",
	srcs = ["test.sh"],
	tags = ["ibazel_live_reload"],
)`)
	e2e.MustWriteFile(t, "test.sh", `
echo "hello moto"`)

	ibazel := e2e.SetUp(t)
	defer ibazel.Kill()
	ibazel.Run([]string{}, "//:live_reload")

	ibazel.ExpectOutput("hello moto")
	out := ibazel.GetOutput()
	t.Logf("Output: '%s'", out)

	if out == "" {
		t.Fatal("Output was empty. Expected at least some output")
	}

	e2e.MustWriteFile(t, "test.sh", `
echo "moto hello"`)
	ibazel.ExpectOutput("moto hello")
	out = ibazel.GetOutput()
	t.Logf("Output: '%s'", out)

	if out == "" {
		t.Fatal("Output was empty. Expected at least some output")
	}
}
