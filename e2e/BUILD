load("@io_bazel_rules_go//go:def.bzl", "go_library")

# gazelle:map_kind go_test go_bazel_test @io_bazel_rules_go//go/tools/bazel_testing:def.bzl

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
        "@io_bazel_rules_go//go/tools/bazel:go_default_library",
        "@io_bazel_rules_go//go/tools/bazel_testing:go_default_library",
    ],
)
