load("@bazel_gazelle//:deps.bzl", "go_repository")

# bazel run //:gazelle -- update-repos -from_file=go.mod -to_macro=repositories.bzl%go_repositories

def go_repositories():
    go_repository(
        name = "com_github_fsnotify_fsnotify",
        importpath = "github.com/fsnotify/fsnotify",
        sum = "h1:hsms1Qyu0jgnwNXIxa+/V/PDsU6CfLf6CNO8H7IWoS4=",
        version = "v1.4.9",
    )
    go_repository(
        name = "com_github_golang_protobuf",
        importpath = "github.com/golang/protobuf",
        sum = "h1:oOuy+ugB+P/kBdUnG5QaMXSIyJ1q38wWSojYCb3z5VQ=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_github_google_go_cmp",
        importpath = "github.com/google/go-cmp",
        sum = "h1:xsAVV57WRhGj6kEIi8ReJzQlHHqcBYCElAvkovg3B/4=",
        version = "v0.4.0",
    )
    go_repository(
        name = "com_github_gorilla_websocket",
        importpath = "github.com/gorilla/websocket",
        sum = "h1:q7AeDBpnBk8AogcD4DSag/Ukw/KV+YhzLj2bP5HvKCM=",
        version = "v1.4.1",
    )
    go_repository(
        name = "com_github_jaschaephraim_lrserver",
        importpath = "github.com/jaschaephraim/lrserver",
        sum = "h1:24NdJ5N6gtrcoeS4JwLMeruKFmg20QdF/5UnX5S/j18=",
        version = "v0.0.0-20171129202958-50d19f603f71",
    )
    go_repository(
        name = "org_golang_x_sys",
        importpath = "golang.org/x/sys",
        sum = "h1:S/FtSvpNLtFBgjTqcKsRpsa6aVsI6iztaz1bQd9BJwE=",
        version = "v0.0.0-20191029155521-f43be2a4598c",
    )
    go_repository(
        name = "org_golang_google_protobuf",
        importpath = "google.golang.org/protobuf",
        sum = "h1:qdOKuR/EIArgaWNjetjgTzgVTAZ+S/WXVrq9HW9zimw=",
        version = "v1.21.0",
    )
    go_repository(
        name = "org_golang_x_xerrors",
        importpath = "golang.org/x/xerrors",
        sum = "h1:E7g+9GITq07hpfrRu66IVDexMakfv52eLZ2CXBWiKr4=",
        version = "v0.0.0-20191204190536-9bdfabe68543",
    )
    go_repository(
        name = "com_github_mattn_go_shellwords",
        importpath = "github.com/mattn/go-shellwords",
        sum = "h1:Y7Xqm8piKOO3v10Thp7Z36h4FYFjt5xB//6XvOrs2Gw=",
        version = "v1.0.10",
    )
