package external_prefixed_package

import (
	"testing"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
)

// This test exercises a package whose name merely *starts with* the string
// "external" (e.g. "//external_api"), as opposed to an actual
// external-repository label such as "//external/foo" or "//external:foo".
//
// See https://github.com/bazelbuild/bazel-watcher/issues/837: labelsToWatch()
// used an overly broad `strings.HasPrefix(label, "//external")` check that
// silently dropped every source file belonging to such a package, so hot
// reload never fired for changes to those files.
const mainFiles = `
-- external_api/BUILD.bazel --
sh_binary(
  name = "external-api",
  srcs = ["main.sh"],
)
-- external_api/main.sh --
printf "Started!"
`

func TestMain(m *testing.M) {
	e2e.TestMain(m, e2e.Args{
		Main: mainFiles,
	})
}

// TestExternalPrefixedPackageSourceIsWatched runs a target that lives in a
// package whose name starts with "external" and then edits one of its source
// files. ibazel must notice the change and rebuild/rerun the target.
//
// Before the fix this hangs: the source label "//external_api:main.sh" is
// filtered out by labelsToWatch(), no watch is ever registered, and the
// second ExpectOutput times out.
func TestExternalPrefixedPackageSourceIsWatched(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//external_api:external-api")
	defer ibazel.Kill()

	ibazel.ExpectOutput("Started!")

	// Edit the source file inside the "external"-prefixed package and make sure
	// ibazel detects the change and reruns the target.
	e2e.MustWriteFile(t, "external_api/main.sh", `printf "Started2!"`)
	ibazel.ExpectOutput("Started2!")
}
