package simple

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
)

const secondaryBuild = `
sh_library(
	name = "lib",
	data = ["lib.sh"],
	visibility = ["//visibility:public"],
)
`

const secondaryLib = `
function say_hello {
	printf "hello!"
}
`

const secondaryLibAlt = `
function say_hello {
	printf "hello2!"
}
`

const mainFiles = `
-- BUILD.bazel --
sh_binary(
	name = "test",
	srcs = ["test.sh"],
	deps = [
		"@secondary//:lib",
	],
)
-- test.sh --
#!/bin/bash
source ../secondary/lib.sh
say_hello
-- WORKSPACE --
local_repository(
    name = "secondary",
    path = "../secondary",
)
`

var (
	secondaryWd  string
	secondaryWd2 string
)

func TestMain(m *testing.M) {
	e2e.TestMain(m, e2e.Args{
		Main: mainFiles,
		SetUp: func() error {
			// Create two secondary workspaces in sibling folders of the main workspace.
			secondaryWd, _ = filepath.Abs(filepath.Join("..", "secondary"))
			secondaryWd2, _ = filepath.Abs(filepath.Join("..", "secondary-2"))

			// Manually create files in the secondary workspaces.
			for _, wd := range []string{secondaryWd, secondaryWd2} {
				if err := os.Mkdir(wd, 0777); err != nil {
					log.Fatalf("os.Mkdir(%q): %v", wd, err)
				}
				for file, contents := range map[string]string{
					".bazelversion": "6.5.0",
					"BUILD.bazel": secondaryBuild,
					"lib.sh":      secondaryLib,
					"WORKSPACE":   "",
					"MODULE.bazel":   "",
				} {
					if err := ioutil.WriteFile(filepath.Join(wd, file), []byte(contents), 0777); err != nil {
						log.Fatalf("Failed to write file %q: %v", file, err)
					}
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

	ibazel.ExpectOutput("hello!", 50 * time.Second)

	ioutil.WriteFile(
		filepath.Join(secondaryWd, "lib.sh"), []byte(secondaryLibAlt), 0777)
	ibazel.ExpectOutput("hello2!")
}

func TestRunWithRepositoryOverrideModifiedFile(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("--override_repository is not implemented in windows")
	}

	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{fmt.Sprintf("--override_repository=secondary=%s", secondaryWd2)}, "//:test")
	defer ibazel.Kill()

	ibazel.ExpectOutput("hello!")

	ioutil.WriteFile(
		filepath.Join(secondaryWd2, "lib.sh"), []byte(secondaryLibAlt), 0777)
	ibazel.ExpectOutput("hello2!")
}
