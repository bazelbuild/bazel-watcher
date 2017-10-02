load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_test(
    name = "go_default_test",
    srcs = ["simple_test.go"],
    data = [
        "//e2e/bazel",
        "//e2e/simple",
        "//ibazel",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/e2e",
    library = ":go_default_library",
    # Run from the root of the repo.
    rundir = ".",
)

go_library(
    name = "go_default_library",
    srcs = [
        "common.go",
        "ibazel.go",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/e2e",
    visibility = ["//visibility:public"],
    deps = ["@com_github_bazelbuild_rules_go//go/tools/bazel:go_default_library"],
)

# I kept running this wrong so I aliased it to make my monkey brain hurt less.
alias(
    name = "e2e",
    actual = ":go_default_test",
)
