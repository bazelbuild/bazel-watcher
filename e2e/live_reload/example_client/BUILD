load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_binary(
    name = "live_reload",
    # Add a simple data dependency that you can modify.
    data = ["test.txt"],
    embed = [":go_default_library"],
    tags = [
        # This tag starts the live_reload server inside iBazel and instructs it to send reload events to webbrowsers.
        "ibazel_live_reload",
        # This tag instructs ibazel to pipe into stdin a event describing actions.
        "ibazel_notify_changes",
    ],
    deps = [
        "//e2e/live_reload/example_client/data:go_default_library",
    ],
)

go_binary(
    name = "no_live_reload",
    data = ["test.txt"],
    embed = [":go_default_library"],
    deps = [
        "//e2e/live_reload/example_client/data:go_default_library",
    ],
)

filegroup(
    name = "files",
    srcs = glob(["*"]),
    data = [
        "//e2e/live_reload/example_client/data:files",
    ],
    visibility = ["//e2e/live_reload:__pkg__"],
)

go_library(
    name = "go_default_library",
    srcs = ["example_client.go"],
    importpath = "github.com/bazelbuild/bazel-watcher/e2e/live_reload/example_client",
    visibility = ["//visibility:private"],
    deps = [
        "//e2e/live_reload/example_client/data:go_default_library",
        "@io_bazel_rules_go//go/tools/bazel:go_default_library",
    ],
)
