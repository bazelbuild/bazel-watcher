load("@io_bazel_rules_go//go/tools/bazel_testing:def.bzl", "go_bazel_test")
load("@io_bazel_rules_webtesting//web:web.bzl", "web_test_suite")

go_bazel_test(
    name = "go_default_test",
    srcs = ["live_reload_test.go"],
    data = [
        "//e2e/live_reload/example_client:files",
        "//ibazel",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/e2e/live_reload",
    tags = [
        # Tag as manual so that it isn't invoked by //.... The test is wrapped
        # by a web_test_suite (below). Please invoke it there.
        #"manual",
    ],
    deps = [
        "//e2e:go_default_library",
        "//e2e/live_reload/example_client/data:go_default_library",
        "@com_github_google_go_cmp//cmp:go_default_library",
        "@com_github_gorilla_websocket//:go_default_library",
        "@com_github_tebeka_selenium//:go_default_library",
        "@io_bazel_rules_go//go/tools/bazel_testing:go_default_library",
        "@io_bazel_rules_webtesting//go/webtest:go_default_library",
    ],
)

web_test_suite(
    name = "web_test",
    args = [
        "-begin_files",
        "$(locations @io_bazel_rules_go//:all_files)",
        "-end_files",
    ],
    browsers = [
        "@io_bazel_rules_webtesting//browsers:chromium-local",
    ],
    data = [
        "//e2e/live_reload/example_client:files",
        "//ibazel",
        # Required by the wrapped test but implicitly added.
        "@io_bazel_rules_go//:all_files",
    ],
    test = ":go_default_test",
)
