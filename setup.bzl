load("@build_bazel_rules_nodejs//:defs.bzl", "yarn_install")
load("@google_bazel_common//:workspace_defs.bzl", "google_common_workspace_rules")
load("@bazel_pandoc//:repositories.bzl", "pandoc_repositories")

def setup():
    yarn_install(
        name = "npm",
        package_json = "@com_github_ubehebe_bazel_runfiles_server//:package.json",
        yarn_lock = "@com_github_ubehebe_bazel_runfiles_server//:yarn.lock",
    )

    google_common_workspace_rules()

    pandoc_repositories()
