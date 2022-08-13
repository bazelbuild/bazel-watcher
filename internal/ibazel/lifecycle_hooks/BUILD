load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "lifecycle_hooks",
    srcs = ["lifecycle_hooks.go"],
    importpath = "github.com/bazelbuild/bazel-watcher/internal/ibazel/lifecycle_hooks",
    visibility = ["//visibility:public"],
    deps = [
        "//internal/ibazel/log",
        "//internal/ibazel/workspace",
        "//third_party/bazel/master/src/main/protobuf/blaze_query",
        "@com_github_mattn_go_shellwords//:go-shellwords",
    ],
)
