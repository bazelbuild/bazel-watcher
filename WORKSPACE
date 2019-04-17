workspace(name = "com_github_ubehebe_bazel_runfiles_server")

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_file")

BAZEL_VERSION = "0.24.0"

git_repository(
    name = "bazel_pandoc",
    remote = "https://github.com/ProdriveTechnologies/bazel-pandoc.git",
    tag = "v0.2",
)

git_repository(
    name = "bazel_skylib",
    remote = "https://github.com/bazelbuild/bazel-skylib.git",
    tag = "0.8.0",
)

git_repository(
    name = "build_bazel_rules_nodejs",
    remote = "https://github.com/bazelbuild/rules_nodejs.git",
    tag = "0.27.8",
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

load("@google_bazel_common//:workspace_defs.bzl", "google_common_workspace_rules")

google_common_workspace_rules()

load("@bazel_pandoc//:repositories.bzl", "pandoc_repositories")

pandoc_repositories()

load("@build_bazel_rules_nodejs//:defs.bzl", "check_bazel_version", "yarn_install")

check_bazel_version(BAZEL_VERSION)

yarn_install(
    name = "npm",
    package_json = "//:package.json",
    yarn_lock = "//:yarn.lock",
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

maven_jar(
    name = "javax_activation",
    artifact = "javax.activation:activation:1.1.1",
)

maven_jar(
    name = "jcommander",
    artifact = "com.beust:jcommander:1.72",
)
