load("@io_bazel_rules_go//go/tools/bazel_testing:def.bzl", "go_bazel_test")

go_bazel_test(
    name = "simple_test",
    srcs = ["simple_test.go"],
    deps = [
        "//internal/e2e",
        "@io_bazel_rules_go//go/tools/bazel_testing",
    ],
)
