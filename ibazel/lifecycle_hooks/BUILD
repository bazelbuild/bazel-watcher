load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["lifecycle_hooks.go"],
    importpath = "github.com/bazelbuild/bazel-watcher/ibazel/lifecycle_hooks",
    visibility = ["//visibility:public"],
    deps = [
        "//ibazel/log:go_default_library",
        "//ibazel/workspace:go_default_library",
        "//third_party/bazel/master/src/main/protobuf/blaze_query:go_default_library",
        "@com_github_mattn_go_shellwords//:go_default_library",
    ],
)
