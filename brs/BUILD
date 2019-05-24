load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_binary(
    name = "main",
    srcs = [
        "RunfilesServer.go",
        "LiveReloadSnippetInjectingHandler.go",
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
        "IntegrationTestRunner.go",
    ],
    visibility = [
        "//visibility:public",
    ],
)
