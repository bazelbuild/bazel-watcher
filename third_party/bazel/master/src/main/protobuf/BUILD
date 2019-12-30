# gazelle:exclude build.pb.go

package(default_visibility = ["//visibility:public"])

load("@rules_proto//proto:defs.bzl", "proto_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

# Apache Version 2.0, January 2004
licenses(["notice"])

proto_library(
    name = "build",
    srcs = ["build.proto"],
)

go_proto_library(
    name = "go_default_library",
    importpath = "github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf",
    proto = ":blaze_query_proto",
)

proto_library(
    name = "blaze_query_proto",
    srcs = ["build.proto"],
)
