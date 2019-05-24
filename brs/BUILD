load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_binary(
    name = "main",
    srcs = [
        "runfiles_server.go",
        "live_reload_snippet_injecting_handler.go",
    ],
    deps = [
        "@com_github_pkg_browser//:go_default_library",
    ],
    visibility = [
        "//visibility:public",
    ],
)

go_binary(
    name = "IntegrationTestRunner",
    srcs = [
        "integration_test_runner.go",
    ],
    visibility = [
        "//visibility:public",
    ],
)
