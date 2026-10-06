package notify_output_groups

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
)

const mainFiles = `
-- BUILD.bazel --
load(":notify_app.bzl", "notify_app")

sh_binary(
    name = "notification_server",
    srcs = ["notification_server.sh"],
)

notify_app(
    name = "app",
    binary = ":notification_server",
    src = "source.txt",
    tags = [
        "ibazel_notify_changes",
        "ibazel_notify_changes_v1",
    ],
)
-- notify_app.bzl --
def _notify_app_impl(ctx):
    extension = ctx.executable.binary.extension
    executable = ctx.actions.declare_file(
        ctx.label.name + ("." + extension if extension else ""),
    )
    ctx.actions.symlink(
        output = executable,
        target_file = ctx.executable.binary,
        is_executable = True,
    )

    generated = ctx.actions.declare_file(ctx.label.name + ".generated.txt")
    ctx.actions.run_shell(
        inputs = [ctx.file.src],
        outputs = [generated],
        arguments = [ctx.file.src.path, generated.path],
        command = "cp $1 $2",
    )

    return [
        DefaultInfo(
            executable = executable,
            runfiles = ctx.attr.binary[DefaultInfo].default_runfiles,
        ),
        OutputGroupInfo(generated = depset([generated])),
    ]

notify_app = rule(
    implementation = _notify_app_impl,
    attrs = {
        "binary": attr.label(executable = True, cfg = "target"),
        "src": attr.label(allow_single_file = True),
    },
    executable = True,
)
-- notification_server.sh --
#!/usr/bin/env bash
printf 'ready\n'
while IFS= read -r line; do
    printf '%s\n' "$line"
done
-- source.txt --
generated output
`

func TestMain(m *testing.M) {
	e2e.TestMain(m, e2e.Args{Main: mainFiles})
}

func TestStructuredNotificationIncludesBazelOutputGroup(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.RunWithAdditionalArgs("//:app", []string{"--notify_output_groups=generated"})
	defer ibazel.Kill()

	ibazel.ExpectOutput("IBAZEL_EVENT ")

	var event struct {
		OutputGroups map[string][]struct {
			Path string `json:"path"`
			URI  string `json:"uri"`
		} `json:"output_groups"`
		OutputGroupsComplete bool `json:"output_groups_complete"`
	}
	for _, line := range strings.Split(ibazel.GetOutput(), "\n") {
		if !strings.HasPrefix(line, "IBAZEL_EVENT ") {
			continue
		}
		if err := json.Unmarshal([]byte(strings.TrimPrefix(line, "IBAZEL_EVENT ")), &event); err != nil {
			t.Fatalf("decode structured notification: %v", err)
		}
		break
	}

	if !event.OutputGroupsComplete {
		t.Fatal("output groups were not marked complete")
	}
	outputs := event.OutputGroups["generated"]
	if len(outputs) != 1 {
		t.Fatalf("generated outputs = %v, want one output", outputs)
	}
	if !strings.HasSuffix(outputs[0].Path, "/app.generated.txt") {
		t.Errorf("generated output path = %q", outputs[0].Path)
	}
	if !strings.HasPrefix(outputs[0].URI, "file://") {
		t.Errorf("generated output URI = %q, want file URI", outputs[0].URI)
	}
}
