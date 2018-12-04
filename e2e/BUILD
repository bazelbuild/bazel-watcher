load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "buffer.go",
        "common.go",
        "ibazel.go",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/e2e",
    visibility = ["//visibility:public"],
    deps = [
        "@com_github_bazelbuild_bazel_integration_testing//go:go_default_library",
        "@com_github_bazelbuild_rules_go//go/tools/bazel:go_default_library",
    ],
)
