load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_file")

# These repos are needed to build //runfiles_server. Downstream bazel repos should call this from
# their WORKSPACE.
def bazel_runfiles_server_repositories():
    git_repository(
        name = "bazel_gazelle",
        remote = "https://github.com/bazelbuild/bazel-gazelle.git",
        tag = "0.17.0",
    )

    git_repository(
        name = "bazel_skylib",
        remote = "https://github.com/bazelbuild/bazel-skylib.git",
        tag = "0.8.0",
    )

    git_repository(
        name = "io_bazel_rules_go",
        remote = "https://github.com/bazelbuild/rules_go.git",
        tag = "0.18.5",
    )

# These repos are only needed for building the examples. Downstream bazel repos  don't need these.
def bazel_runfiles_server_example_repositories():
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

