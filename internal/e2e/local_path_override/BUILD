load("@io_bazel_rules_go//go/tools/bazel_testing:def.bzl", "go_bazel_test")

go_bazel_test(
    name = "local_path_override_test",
    srcs = ["local_path_override_test.go"],
    deps = [
        "//internal/e2e",
        "@io_bazel_rules_go//go/tools/bazel_testing:go_default_library",
    ],
)
