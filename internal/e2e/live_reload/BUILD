load("@io_bazel_rules_go//go/tools/bazel_testing:def.bzl", "go_bazel_test")

go_bazel_test(
    name = "live_reload_test",
    srcs = ["live_reload_test.go"],
    deps = [
        "//internal/e2e",
        "@com_github_gorilla_websocket//:websocket",
    ],
)
