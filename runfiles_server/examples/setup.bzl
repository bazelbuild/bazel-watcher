load("@build_bazel_rules_nodejs//:defs.bzl", "yarn_install")
load("@bazel_pandoc//:repositories.bzl", "pandoc_repositories")

def setup_examples():
    yarn_install(
        name = "npm",
        package_json = "//runfiles_server/examples:package.json",
        yarn_lock = "//runfiles_server/examples:yarn.lock",
    )

    pandoc_repositories()