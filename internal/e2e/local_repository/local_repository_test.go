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

const mainWorkspaceTailDep = `
function say_local_tail {
	printf "local-tail-1!"
}
`

const mainWorkspaceTailDepAlt = `
function say_local_tail {
	printf "local-tail-2!"
}
`

const mainFiles = `
-- BUILD.bazel --
sh_library(
	name = "late2",
	srcs = ["late2.sh"],
	visibility = ["//visibility:public"],
)

sh_library(
	name = "late",
	srcs = ["late.sh"],
	deps = [":late2"],
	visibility = ["//visibility:public"],
)

sh_binary(
	name = "test",
	srcs = ["test.sh"],
	deps = [
		"@secondary//:lib",
		":late",
	],
)
-- test.sh --
#!/bin/bash
source ../secondary/lib.sh
source late.sh
source late2.sh
say_hello
say_local_tail
-- late.sh --
function say_local {
	true
}
-- late2.sh --
function say_local_tail {
	printf "local-tail-1!"
}
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

func TestRunWithModifiedWorkspaceFileAfterLocalRepositoryDependency(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skipf("local_repository is not implemented in windows")
	}

	ibazel := e2e.SetUp(t)
	ioutil.WriteFile("late2.sh", []byte(mainWorkspaceTailDep), 0777)
	ibazel.Run([]string{}, "//:test")
	defer ibazel.Kill()

	// Reproduces the bug fixed in 9fbb5bc: cquery can interleave labels as
	// local and workspace labels, so workspace files after @repo labels must
	// still remain watched.
	ibazel.ExpectOutput("local-tail-1!", 50*time.Second)

	ioutil.WriteFile("late2.sh", []byte(mainWorkspaceTailDepAlt), 0777)
	ibazel.ExpectOutput("local-tail-2!")
}
