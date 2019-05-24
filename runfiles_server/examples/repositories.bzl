load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

# These repos are only needed for building the examples. Downstream bazel repos don't need these.
def example_repositories():
    git_repository(
        name = "bazel_pandoc",
        remote = "https://github.com/ProdriveTechnologies/bazel-pandoc.git",
        tag = "v0.2",
    )

    git_repository(
        name = "build_bazel_rules_nodejs",
        remote = "https://github.com/bazelbuild/rules_nodejs.git",
        tag = "0.27.8",
    )

