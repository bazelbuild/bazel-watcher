load("@io_bazel_rules_go//go:def.bzl", "go_library")

# gazelle:map_kind go_test go_bazel_test @io_bazel_rules_go//go/tools/bazel_testing:def.bzl

go_library(
    name = "e2e",
    srcs = [
        "buffer.go",
        "common.go",
        "ibazel.go",
    ],
    data = [
        "//cmd/ibazel",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/internal/e2e",
    visibility = ["//visibility:public"],
    deps = [
        "@io_bazel_rules_go//go/tools/bazel",
        "@io_bazel_rules_go//go/tools/bazel_testing",
        "@org_golang_x_tools//txtar",
    ],
)
