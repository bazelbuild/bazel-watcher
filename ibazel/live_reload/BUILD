load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "events.go",
        "server.go",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/ibazel/live_reload",
    visibility = ["//visibility:public"],
    deps = [
        "//ibazel/log:go_default_library",
        "//third_party/bazel/master/src/main/protobuf/blaze_query:go_default_library",
        "@com_github_jaschaephraim_lrserver//:go_default_library",
    ],
)
