package termination

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/e2e"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

const mainFiles = `
-- sleep.bzl --
def _sleep(ctx):
    inputs = ctx.files.srcs[:]
    out = ctx.actions.declare_file("out")
    ctx.actions.run_shell(
        command = "sleep 8; touch {}".format(out.path),
        inputs = inputs,
        outputs = [out],
        execution_requirements = {
            "timeout": "1",
        },
    )
    return [DefaultInfo(files = depset([out]))]

sleep = rule(
    implementation = _sleep,
    attrs = {
        "srcs": attr.label_list(default = [], allow_files = True),
	}
)

-- BUILD.bazel --
load("sleep.bzl", "sleep")

sleep(
	name = "abort_compilation_early",
	srcs = [
		"//:abort_compilation_early.sh",
	]
)
-- abort_compilation_early.sh --
printf "Started!"
`

const buildSuccStr = `Build completed successfully`
const buildChangedStr = `Changed:`
const buildCancelStr = `Cancelling previous Bazel invocation and rebuilding...`

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: mainFiles,
	})
}

func TestAbortCompilationEarlyWithoutChanges(t *testing.T) {
	os.Setenv("IBAZEL_ABORT_COMPILATION_EARLY", "1")
	ibazel := e2e.SetUp(t)
	defer ibazel.Kill()

	// wait for a complete start and assert we have a build success log
	ibazel.Build("//:abort_compilation_early")
	time.Sleep(10 * 1000 * time.Millisecond)
	ibazel.ExpectError(buildSuccStr)

	defer ibazel.Kill()
}

func TestAbortCompilationEarlyAfterChange(t *testing.T) {
	os.Setenv("IBAZEL_ABORT_COMPILATION_EARLY", "1")
	ibazel := e2e.SetUp(t)
	defer ibazel.Kill()

	// wait for a complete start and assert we have a build success log
	ibazel.Build("//:abort_compilation_early")
	time.Sleep(10 * 1000 * time.Millisecond)
	ibazel.ExpectError(buildSuccStr)
	
	// assert we have a first 'changed:' log
	e2e.MustWriteFile(t, "abort_compilation_early.sh", "printf \"Started1!\";")
	time.Sleep(4 * 1000 * time.Millisecond)
	ibazel.ExpectError(buildChangedStr)

	// assert we have a first 'cancelling previous Bazel invocation and rebuilding...' log
	e2e.MustWriteFile(t, "abort_compilation_early.sh", "printf \"Started2!\";")
	time.Sleep(2 * 1000 * time.Millisecond)
	ibazel.ExpectError(buildCancelStr)

	// push two additional abort triggers
	e2e.MustWriteFile(t, "abort_compilation_early.sh", "printf \"Started3!\";")
	time.Sleep(2 * 1000 * time.Millisecond)
	ibazel.ExpectError(buildCancelStr)

	e2e.MustWriteFile(t, "abort_compilation_early.sh", "printf \"Started4!\";")
	time.Sleep(2 * 1000 * time.Millisecond)
	ibazel.ExpectError(buildCancelStr)

	time.Sleep(10 * 1000 * time.Millisecond)
	out := ibazel.GetError()

	// assert we have a build success log
	number_of_cancelled_invocations := strings.Count(out, buildCancelStr)
	defer ibazel.Kill()
	if number_of_cancelled_invocations != 3 {
		t.Errorf("Expected number of cancelled invokations was 3 but found %d", number_of_cancelled_invocations)
	}
}
