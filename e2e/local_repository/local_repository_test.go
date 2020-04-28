package simple

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-watcher/e2e"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
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
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: mainFiles,
		SetUp: func() error {
			// Create two secondary workspaces in sibling folders of the main workspace.
			secondaryWd, _ = filepath.Abs(filepath.Join("..", "secondary"))
			secondaryWd2, _ = filepath.Abs(filepath.Join("..", "secondary-2"))

			// Manually create files in the secondary workspaces.
			for _, wd := range []string{secondaryWd, secondaryWd2} {
				os.Mkdir(wd, 0777)
				ioutil.WriteFile(
					filepath.Join(wd, "BUILD.bazel"), []byte(secondaryBuild), 0777)
				ioutil.WriteFile(
					filepath.Join(wd, "lib.sh"), []byte(secondaryLib), 0777)
				ioutil.WriteFile(
					filepath.Join(wd, "WORKSPACE"), []byte(""), 0777)
			}
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			if err := filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if strings.HasSuffix(path, ".sh") {
					if err := os.Chmod(path, 0777); err != nil {
						return fmt.Errorf("Error os.Chmod(%q, 0777): %v", path, err)
					}
				}
				return nil
			}); err != nil {
				fmt.Printf("Error walking dir: %v\n", err)
				return err
			}
			return nil
		},
	})
}

func TestRunWithModifiedFile(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//:test")
	defer ibazel.Kill()

	ibazel.ExpectOutput("hello!")

	ioutil.WriteFile(
		filepath.Join(secondaryWd, "lib.sh"), []byte(secondaryLibAlt), 0777)
	ibazel.ExpectOutput("hello2!")
}

func TestRunWithRepositoryOverrideModifiedFile(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{fmt.Sprintf("--override_repository=secondary=%s", secondaryWd2)}, "//:test")
	defer ibazel.Kill()

	ibazel.ExpectOutput("hello!")

	ioutil.WriteFile(
		filepath.Join(secondaryWd2, "lib.sh"), []byte(secondaryLibAlt), 0777)
	ibazel.ExpectOutput("hello2!")
}
