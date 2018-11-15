# gazelle:exclude profiler_test.go

load("@com_github_bazelbuild_bazel_integration_testing//go:bazel_integration_test.bzl", "bazel_go_integration_test")

bazel_go_integration_test(
    name = "go_default_test",
    srcs = ["profiler_test.go"],
    data = [
        "//ibazel",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/e2e/profiler",
    deps = [
        "//e2e:go_default_library",
        "@com_github_bazelbuild_bazel_integration_testing//go:go_default_library",
        "@com_github_gorilla_websocket//:go_default_library",
    ],
)
