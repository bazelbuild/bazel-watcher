load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_binary(
    name = "live_reload",
    # Add a simple data dependency that you can modify.
    data = ["data.txt"],
    embed = [":example_client_lib"],
    tags = [
        # This tag starts the live_reload server inside iBazel and instructs it to send reload events to webbrowsers.
        "ibazel_live_reload",
        # This tag instructs ibazel to pipe into stdin a event describing actions.
        "ibazel_notify_changes",
    ],
    deps = [
        "//example_client/data",
    ],
)

go_binary(
    name = "no_live_reload",
    data = ["data.txt"],
    embed = [":example_client_lib"],
    deps = [
        "//example_client/data",
    ],
)

filegroup(
    name = "files",
    srcs = glob(["*"]),
    data = [
        "//example_client/data:files",
    ],
    visibility = ["//internal/e2e/example_client:__pkg__"],
)

go_library(
    name = "example_client_lib",
    srcs = ["example_client.go"],
    importpath = "github.com/bazelbuild/bazel-watcher/example_client",
    visibility = ["//visibility:private"],
    deps = [
        "//example_client/data",
        "@io_bazel_rules_go//go/tools/bazel",
    ],
)
