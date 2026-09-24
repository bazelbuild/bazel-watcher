load("@io_bazel_rules_go//go/tools/bazel_testing:def.bzl", "go_bazel_test")

go_bazel_test(
    name = "external_prefixed_package_test",
    srcs = ["external_prefixed_package_test.go"],
    deps = [
        "//internal/e2e",
        "@io_bazel_rules_go//go/tools/bazel_testing:go_default_library",
    ],
)
