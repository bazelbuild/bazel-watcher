# gazelle:exclude simple_test.go

"""TODO(achew): write a module docstring"""

load("@build_bazel_integration_testing//go:bazel_integration_test.bzl", "bazel_go_integration_test")

bazel_go_integration_test(
    name = "go_default_test",
    srcs = ["simple_test.go"],
    data = [
        "//ibazel",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/e2e/simple",
    deps = [
        "//e2e:go_default_library",
        "@build_bazel_integration_testing//go:go_default_library",
    ],
)
