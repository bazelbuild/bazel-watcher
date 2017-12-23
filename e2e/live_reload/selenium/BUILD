# gazelle:exclude live_reload_test.go

load("@io_bazel_rules_go//go:def.bzl", "go_embed_data")
load("@io_bazel_rules_webtesting//web/internal:wrap_web_test_suite.bzl", "web_test_suite")
load("@com_github_bazelbuild_bazel_integration_testing//go:bazel_integration_test.bzl", "bazel_go_integration_test")
load("@com_github_bazelbuild_bazel_integration_testing//tools:common.bzl", "BAZEL_VERSIONS")

bazel_go_integration_test(
    name = "bazel_integration",
    srcs = [
        "live_reload_test.go",
        ":python",
    ],
    data = [
        "//ibazel",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/e2e/live_reload/selenium",
    # Since these tests are used by the web_test_suite below, don't run them in
    # //... queries.
    tags = ["manual"],
    deps = [
        "//e2e:go_default_library",
        "@com_github_bazelbuild_bazel_integration_testing//go:go_default_library",
        "@com_github_bazelbuild_rules_webtesting//go/webtest:go_default_library",
        "@com_github_gorilla_websocket//:go_default_library",
        "@com_github_tebeka_selenium//:go_default_library",
    ],
)

[web_test_suite(
    name = "go_default_test/bazel%s" % v,
    browsers = [
        "@io_bazel_rules_webtesting//browsers:chromium-local",
    ],
    tags = {
        "default": ["notravis"],
    },
    test = ":bazel_integration/bazel0.16.1",
) for v in BAZEL_VERSIONS]

go_embed_data(
    name = "python",
    src = "live_reload.py",
    var = "python",
)

py_binary(
    name = "live_reload",
    srcs = ["live_reload.py"],
    data = ["test.txt"],
    tags = [
        "ibazel_live_reload",
        "ibazel_notify_changes",
    ],
)
