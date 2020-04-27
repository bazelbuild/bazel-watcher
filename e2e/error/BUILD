load("@io_bazel_rules_go//go/tools/bazel_testing:def.bzl", "go_bazel_test")

go_bazel_test(
    name = "go_default_test",
    srcs = ["error_test.go"],
    data = [
        "//ibazel",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/e2e/error",
    deps = [
        "//e2e:go_default_library",
        "//e2e/example_client:go_default_library",
    ],
)
