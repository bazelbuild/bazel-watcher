package error

import (
	"testing"

	"github.com/bazelbuild/bazel-watcher/e2e"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{})
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
