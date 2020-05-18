load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["data.go"],
    importpath = "github.com/bazelbuild/bazel-watcher/example_client/data",
    visibility = ["//visibility:public"],
)

filegroup(
    name = "files",
    srcs = glob(["*"]),
    visibility = ["//example_client:__pkg__"],
)
