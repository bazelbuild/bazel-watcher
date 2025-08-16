package windows_shell_smoke

import (
	"runtime"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
)

const mainFiles = `
-- MODULE.bazel --
bazel_dep(name = "rules_shell", version = "0.2.0")

-- BUILD.bazel --
load("@rules_shell//shell:sh_binary.bzl", "sh_binary")

sh_binary(
  name = "smoke",
  srcs = ["smoke.sh"],
)

-- smoke.sh --
#!/usr/bin/env bash
printf "shell-ok"
`

func TestMain(m *testing.M) {
	e2e.TestMain(m, e2e.Args{Main: mainFiles})
}

func TestWindowsShellPresent(t *testing.T) {
	// Only meaningful on Windows — skip elsewhere.
	if runtime.GOOS != "windows" {
		t.Skip("windows-only smoke test")
	}
	ibazel := e2e.SetUp(t)
	defer ibazel.Kill()
	ibazel.Run([]string{}, "//:smoke")
	ibazel.ExpectOutput("shell-ok", 35*time.Second)
}
