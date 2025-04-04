package simple

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
)

const mainFiles = `
-- BUILD.bazel --
constraint_setting(
	name = "constraint_setting",
	default_constraint_value = ":constraint1",
)
constraint_value(
	name = "constraint1",
	constraint_setting = "constraint_setting",
)
constraint_value(
	name = "constraint2",
	constraint_setting = "constraint_setting",
)

platform(
	name = "platform1",
	constraint_values = [
		":constraint1",
	],
)
platform(
	name = "platform2",
	constraint_values = [
		":constraint2",
	],
)

sh_binary(
	name = "incompatible_by_default",
	srcs = ["incompatible_by_default.sh"],
	target_compatible_with = [
		":constraint2",
	],
)
-- incompatible_by_default.sh --
#!/bin/bash
echo 'hello!'
`

var secondaryWd string

func TestMain(m *testing.M) {
	e2e.TestMain(m, e2e.Args{
		Main: mainFiles,
		SetUp: func() error {
			path := "incompatible_by_default.sh"
			if err := os.Chmod(path, 0777); err != nil {
				return fmt.Errorf("Error os.Chmod(%q, 0777): %v", path, err)
			}
			return nil
		},
	})
}

func TestRunWithPlatforms(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{"--platforms=//:platform2"}, "//:incompatible_by_default")
	defer ibazel.Kill()

	ibazel.ExpectOutput("hello!", 35 * time.Second)
}
