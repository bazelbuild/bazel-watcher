workspace(name = "com_github_ubehebe_bazel_runfiles_server")

BAZEL_VERSION = "0.24.0"

load(":repositories.bzl", "bazel_runfiles_server_repositories")

bazel_runfiles_server_repositories()

load("@build_bazel_rules_nodejs//:defs.bzl", "check_bazel_version")

check_bazel_version(BAZEL_VERSION)

load(":setup.bzl", "setup")

setup()
