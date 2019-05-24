load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

def setup():
    go_repository(
        name = "com_github_pkg_browser",
        importpath = "github.com/pkg/browser",
        commit = "0a3d74bf9ce488f035cf5bc36f753a711bc74334",
    )

    go_rules_dependencies()
    go_register_toolchains()
    gazelle_dependencies()
