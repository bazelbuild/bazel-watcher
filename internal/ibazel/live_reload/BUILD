load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "live_reload",
    srcs = [
        "events.go",
        "server.go",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/internal/ibazel/live_reload",
    visibility = ["//visibility:public"],
    deps = [
        "//internal/ibazel/log",
        "//third_party/bazel/master/src/main/protobuf/blaze_query",
        "@com_github_jaschaephraim_lrserver//:lrserver",
    ],
)
