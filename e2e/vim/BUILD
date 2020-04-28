load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//go/tools/bazel_testing:def.bzl", "go_bazel_test")

go_bazel_test(
    name = "go_default_test",
    srcs = ["vim.go"],
    data = [
        "//ibazel",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/e2e/vim",
    deps = [
        "//bazel:go_default_library",
        "//e2e:go_default_library",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["vim.go"],
    importpath = "github.com/bazelbuild/bazel-watcher/e2e/vim",
    visibility = ["//visibility:public"],
    deps = [
        "//e2e:go_default_library",
        "@io_bazel_rules_go//go/tools/bazel_testing:go_default_library",
    ],
)
