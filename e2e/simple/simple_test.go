// Copyright 2017 The Bazel Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package simple

import (
	"os"
	"runtime/debug"
	"testing"

	bazel "github.com/bazelbuild/bazel-integration-testing/go"
	"github.com/bazelbuild/bazel-watcher/e2e"
)

func must(t *testing.T, e error) {
	if e != nil {
		t.Errorf("Error: %s", e)
		debug.PrintStack()
	}
}

func TestSimpleRun(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchFileWithMode("test.sh", `printf "Started!"`, 0777))
	must(t, b.ScratchFile("BUILD", `
sh_binary(
	name = "test",
	srcs = ["test.sh"],
)
`))

	ibazel := e2e.NewIBazelTester(t, b)
	ibazel.Run("//:test")
	defer ibazel.Kill()

	ibazel.ExpectOutput("Started!")
}

func TestSimpleRunUnderSubdir(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchDir("subdir"))
	must(t, b.ScratchFileWithMode("subdir/test.sh", `printf "Started!"`, 0777))
	must(t, b.ScratchFile("subdir/BUILD", `
sh_binary(
	name = "test",
	srcs = ["test.sh"],
)
`))

	ibazel := e2e.NewIBazelTester(t, b)

	err = os.Chdir("subdir")
	if err != nil {
		t.Fatal(err)
	}

	ibazel.Run("test")
	defer ibazel.Kill()

	err = os.Chdir("..")
	if err != nil {
		t.Fatal(err)
	}

	ibazel.ExpectOutput("Started!")
}

func TestSimpleRunWithModifiedFile(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchFileWithMode("test.sh", `printf "Started!"`, 0777))
	must(t, b.ScratchFile("BUILD", `
sh_binary(
	name = "test",
	srcs = ["test.sh"],
)
`))

	ibazel := e2e.NewIBazelTester(t, b)
	ibazel.Run("//:test")
	defer ibazel.Kill()

	// Give it time to start up and query.
	must(t, b.ScratchFileWithMode("test.sh", `printf "Started2!"`, 0777))
	ibazel.ExpectOutput("Started2!")

	// Manipulate a source file and sleep past the debounce.
	must(t, b.ScratchFileWithMode("test.sh", `printf "Started3!"`, 0777))
	ibazel.ExpectOutput("Started3!")

	// Now a BUILD file.
	must(t, b.ScratchFile("BUILD", `
sh_binary(
	# New comment
	name = "test",
	srcs = ["test.sh"],
)
`))
	ibazel.ExpectOutput("Started3!")
}
