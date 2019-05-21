load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_file")

# These repos are needed to build //brs. Downstream bazel repos should call this from their
# WORKSPACE.
def bazel_runfiles_server_repositories():
    git_repository(
        name = "bazel_skylib",
        remote = "https://github.com/bazelbuild/bazel-skylib.git",
        tag = "0.8.0",
    )

    git_repository(
        name = "flogger",
        remote = "https://github.com/google/flogger.git",
        tag = "flogger-0.2",
    )

    git_repository(
        name = "google_bazel_common",
        commit = "758c17dbc7b724e64f915b5496708cd01ffd38d5",
        remote = "https://github.com/google/bazel-common.git",
    )

    git_repository(
        name = "io_bazel_rules_go",
        remote = "https://github.com/bazelbuild/rules_go.git",
        tag = "0.18.5",
    )

    # The official mime.types file from apache httpd. This is way more up-to-date than the default jdk one
    # (javax.activation).
    http_file(
        name = "apache_mime_types",
        sha256 = "ce95e59e1f7fed0ebda42b61de49d213f85ed31b73214cccf08b47cb98f81814",
        urls = [
            # TODO: use 2.5 when it's released. 2.4 doesn't have font/woff2 (used by katex)
            "https://svn.apache.org/repos/asf/httpd/httpd/tags/2.5.0-alpha/docs/conf/mime.types",
        ],
    )

    native.maven_jar(
        name = "javax_activation",
        artifact = "javax.activation:activation:1.1.1",
    )

    native.maven_jar(
        name = "jcommander",
        artifact = "com.beust:jcommander:1.72",
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

