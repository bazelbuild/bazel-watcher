load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_binary(
    name = "live_reload",
    # Add a simple data dependency that you can modify.
    data = ["data.txt"],
    embed = [":go_default_library"],
    tags = [
        # This tag starts the live_reload server inside iBazel and instructs it to send reload events to webbrowsers.
        "ibazel_live_reload",
        # This tag instructs ibazel to pipe into stdin a event describing actions.
        "ibazel_notify_changes",
    ],
    deps = [
        "//example_client/data:go_default_library",
    ],
)

go_binary(
    name = "no_live_reload",
    data = ["data.txt"],
    embed = [":go_default_library"],
    deps = [
        "//example_client/data:go_default_library",
    ],
)

filegroup(
    name = "files",
    srcs = glob(["*"]),
    data = [
        "//example_client/data:files",
    ],
    visibility = ["//e2e/example_client:__pkg__"],
)

go_library(
    name = "go_default_library",
    srcs = ["example_client.go"],
    importpath = "github.com/bazelbuild/bazel-watcher/example_client",
    visibility = ["//visibility:private"],
    deps = [
        "//example_client/data:go_default_library",
        "@io_bazel_rules_go//go/tools/bazel:go_default_library",
    ],
)
