"""TODO(achew): write a module docstring"""

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "common.go",
        "ibazel.go",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/e2e",
    visibility = ["//visibility:public"],
    deps = [
        "@build_bazel_integration_testing//go:go_default_library",
        "@com_github_bazelbuild_rules_go//go/tools/bazel:go_default_library",
    ],
)
