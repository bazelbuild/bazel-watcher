workspace(name = "com_github_ubehebe_bazel_runfiles_server")

BAZEL_VERSION = "0.24.0"

load(":repositories.bzl", "bazel_runfiles_server_repositories", "bazel_runfiles_server_example_repositories")

bazel_runfiles_server_repositories()

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

bazel_runfiles_server_example_repositories()

load("@bazel_skylib//lib:versions.bzl", "versions")

versions.check(minimum_bazel_version=BAZEL_VERSION)

load(":setup.bzl", "setup")

setup()

load("@build_bazel_rules_nodejs//:defs.bzl", "yarn_install")

yarn_install(
    name = "npm",
    package_json = "@com_github_ubehebe_bazel_runfiles_server//:package.json",
    yarn_lock = "@com_github_ubehebe_bazel_runfiles_server//:yarn.lock",
)

load("@bazel_pandoc//:repositories.bzl", "pandoc_repositories")

pandoc_repositories()

