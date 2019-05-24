load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")
load("//runfiles_server:serve.bzl", "serve")
load("@bazel_pandoc//:pandoc.bzl", "pandoc")

go_binary(
    name = "runfiles_server",
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



# A nice livereload experience for writing the readme itself.

pandoc(
    name = "run_pandoc",
    src = "README.md",
    from_format = "markdown",
    to_format = "html",
)

serve(
    name = "readme",
    index = ":run_pandoc",
)
