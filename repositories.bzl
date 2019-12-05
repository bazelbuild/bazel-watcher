load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_repositories():
    go_repository(
        name = "com_github_fsnotify_fsnotify",
        importpath = "github.com/fsnotify/fsnotify",
        tag = "v1.4.7",
    )

    go_repository(
        name = "com_github_jaschaephraim_lrserver",
        importpath = "github.com/jaschaephraim/lrserver",
        tag = "3.0.1",
    )

    go_repository(
        name = "com_github_gorilla_websocket",
        importpath = "github.com/gorilla/websocket",
        tag = "v1.4.1",
    )

    go_repository(
        name = "org_golang_x_sys",
        commit = "cc5685c2db1239775905f3911f0067c0fa74762f",
        importpath = "golang.org/x/sys",
    )

    # NOTE: this must match rules_go version located at
    # https://github.com/bazelbuild/rules_go/blob/master/go/private/repositories.bzl
    go_repository(
        name = "com_github_golang_protobuf",
        importpath = "github.com/golang/protobuf",
        tag = "v1.3.2",
    )

    go_repository(
        name = "com_github_gorilla_websocket",
        importpath = "github.com/gorilla/websocket",
        tag = "v1.4.0",
    )

    # NOTE: this must match rules_go version above
    go_repository(
        name = "com_github_bazelbuild_rules_go",
        importpath = "github.com/bazelbuild/rules_go",
        tag = "v0.20.3",
    )
