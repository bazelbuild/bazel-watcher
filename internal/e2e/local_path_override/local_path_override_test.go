package simple

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
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

const secondaryModule = `
module(name = "secondary")
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
# Due to a lack of runfiles file support, search any potential override path.
for f in "../secondary*/lib.sh"; do source $f; done
say_hello

-- MODULE.bazel --
bazel_dep(name = "secondary")
local_path_override(
    module_name = "secondary",
    path = "../local_secondary",
)
`

var (
	secondaryWd string
)

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: mainFiles,
		SetUp: func() error {
			// Create a secondary repository in a sibling folder.
			secondaryWd, _ = filepath.Abs(filepath.Join("..", "local_secondary"))

			// Manually create files in the secondary repository.
			if err := os.Mkdir(secondaryWd, 0777); err != nil {
				log.Fatalf("os.Mkdir(%q): %v", secondaryWd, err)
			}
			for file, contents := range map[string]string{
				"BUILD.bazel": secondaryBuild,
				"lib.sh":      secondaryLib,
				"MODULE.bazel":   secondaryModule,
				"WORKSPACE":   "",
			} {
				if err := ioutil.WriteFile(filepath.Join(secondaryWd, file), []byte(contents), 0777); err != nil {
					log.Fatalf("Failed to write file %q: %v", file, err)
				}
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

func TestRunWithOverrideRepository(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{"--enable_bzlmod=1"}, "//:test")
	defer ibazel.Kill()

	ibazel.ExpectOutput("hello!", 35 * time.Second)

	if err := ioutil.WriteFile(filepath.Join(secondaryWd, "lib.sh"), []byte(secondaryRecompiledLib), 0777); err != nil {
		log.Fatalf("Failed to write file lib.sh a second time. (%v)", err)
	}
	ibazel.ExpectOutput("hello!", 35 * time.Second)
}

