package simple

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
)

const nestedBuild = `
sh_library(
	name = "lib",
	data = ["lib.sh"],
	visibility = ["//visibility:public"],
)
`

const nestedLib = `
function say_hello {
	printf "hello!"
}
`

const nestedLibAlt = `
function say_hello {
	printf "hello2!"
}
`

const mainFiles = `
-- .bazelversion --
6.5.0
-- BUILD.bazel --
sh_binary(
    name = "test",
    srcs = ["test.sh"],
    args = ["$(location @nested//:lib)"],
    deps = ["@nested//:lib"],
)
-- test.sh --
#!/bin/bash
source "$1"
say_hello
`

const mainModuleFile = `
module(name = "primary")
bazel_dep(name = "nested_module", repo_name = "nested")
local_path_override(
    module_name = "nested_module",
    path = "./nested/",
)
`

const nestedModuleFile = `
module(
	name = "nested_module",
	repo_name = "nested",
)
`

var nestedWd string

func TestMain(m *testing.M) {
	e2e.TestMain(m, e2e.Args{
		Main:              mainFiles,
		ModuleFileContent: mainModuleFile,
		SetUp: func() error {
			// Create a nested module in a subfolder.
			nestedWd, _ = filepath.Abs("nested")

			// Manually create files in the nested module.
			if err := os.Mkdir(nestedWd, 0777); err != nil {
				log.Fatalf("os.Mkdir(%q): %v", nestedWd, err)
			}
			for file, contents := range map[string]string{
				".bazelversion": "6.5.0", // Needed for built-in sh_binary.
				"BUILD.bazel":   nestedBuild,
				"lib.sh":        nestedLib,
				"WORKSPACE":     "# No content since we are using MODULE.bazel.",
				"MODULE.bazel":  nestedModuleFile,
			} {
				if err := ioutil.WriteFile(filepath.Join(nestedWd, file), []byte(contents), 0777); err != nil {
					log.Fatalf("Failed to write file %q: %v", file, err)
				}
			}

			return nil
		},
	})
}

func TestRunWithModifiedFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("--override_repository is not implemented in windows")
	}

	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//:test")
	defer ibazel.Kill()

	ibazel.ExpectOutput("hello!")

	ioutil.WriteFile(
		filepath.Join(nestedWd, "lib.sh"), []byte(nestedLibAlt), 0777)
	ibazel.ExpectOutput("hello2!")
}
