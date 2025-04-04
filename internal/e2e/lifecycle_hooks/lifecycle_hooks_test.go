package lifecycle_hooks

import (
	"testing"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
)

const mainFiles = `
-- BUILD.bazel --
sh_binary(
  name = "test",
  srcs = ["test.sh"],
)
sh_binary(
  name = "failure",
  srcs = ["failure.sh"],
)
-- test.sh --
printf "action"
`

func TestMain(m *testing.M) {
	e2e.TestMain(m, e2e.Args{
		Main: mainFiles,
	})
}

func TestLifecycleHooks(t *testing.T) {
	ibazel := e2e.SetUp(t)
	defer ibazel.Kill()

	ibazel.RunWithAdditionalArgs("//:test", []string{
		"-run_command_before=echo hi-before",
		"-run_command_after=echo hi-after",
		"-run_command_after_success=echo hi-after-success",
		"-run_command_after_error=echo hi-after-error",
	})
	ibazel.ExpectOutput("hi-before")
	ibazel.ExpectOutput("hi-after")
	ibazel.ExpectOutput("hi-after-success")
}

func TestLifecycleHooksFailure(t *testing.T) {
	ibazel := e2e.SetUp(t)
	defer ibazel.Kill()

	ibazel.RunUnverifiedWithAdditionalArgs("//:failure", []string{
		"-run_command_before=echo hi-before",
		"-run_command_after=echo hi-after",
		"-run_command_after_success=echo hi-after-success",
		"-run_command_after_error=echo hi-after-error",
	})
	ibazel.ExpectOutput("hi-before")
	ibazel.ExpectOutput("hi-after")
	ibazel.ExpectOutput("hi-after-error")
}
