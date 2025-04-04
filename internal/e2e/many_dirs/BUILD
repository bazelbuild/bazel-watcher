load("@io_bazel_rules_go//go/tools/bazel_testing:def.bzl", "go_bazel_test")

go_bazel_test(
    name = "many_dirs",
    srcs = ["many_dirs_test.go"],
    deps = [
        "//internal/e2e",
        "@io_bazel_rules_go//go/tools/bazel_testing:go_default_library",
    ],
)

go_bazel_test(
    name = "many_dirs_test",
    srcs = ["many_dirs_test.go"],
    deps = ["//internal/e2e"],
)
