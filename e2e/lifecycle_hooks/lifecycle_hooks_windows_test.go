package lifecycle_hooks

import (
	"testing"

	"github.com/bazelbuild/bazel-watcher/e2e"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

const mainFiles = `
-- BUILD.bazel --
sh_binary(
  name = "test",
  srcs = ["test.sh"],
)
-- test.sh --
printf "action"

-- command_before.sh --
#!/bin/sh
printf "Hello from script"
`

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: mainFiles,
	})
}

func TestLifecycleHooks(t *testing.T) {
	ibazel := e2e.SetUp(t)

	ibazel.RunWithAdditionalArgs("//:test", []string{
		`-run_command_before=echo Hello`,
	})
	ibazel.ExpectOutput("Hello")
	ibazel.Kill()
}
