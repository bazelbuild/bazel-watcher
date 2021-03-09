load("@io_bazel_rules_go//go/tools/bazel_testing:def.bzl", "go_bazel_test")

go_bazel_test(
    name = "go_default_test",
    srcs = ["notify_changes_test.go"],
    data = [
        "//ibazel",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/e2e/notify_changes",
    deps = [
        "//e2e:go_default_library",
        "@com_github_gorilla_websocket//:go_default_library",
        "@io_bazel_rules_go//go/tools/bazel_testing:go_default_library",
    ],
)
