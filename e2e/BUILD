# gazelle:exclude simple_test.go

load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@com_github_bazelbuild_bazel_integration_testing//go:bazel_integration_test.bzl", "bazel_go_integration_test")

bazel_go_integration_test(
    name = "go_default_test",
    srcs = ["simple_test.go"],
    data = [
        "//ibazel",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/e2e",
    library = ":go_default_library",
)

go_library(
    name = "go_default_library",
    srcs = [
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

# I kept running this wrong so I aliased it to make my monkey brain hurt less.
# Deactivated until
# https://github.com/bazelbuild/bazel/issues/4205
#alias(
#    name = "e2e",
#    actual = "//e2e:go_default_test/bazel0.7",
#)
