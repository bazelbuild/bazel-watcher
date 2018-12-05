# Copyright 2017 The Bazel Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_binary(
    name = "ibazel",
    embed = [":go_default_library"],
    importpath = "github.com/bazelbuild/bazel-watcher/ibazel",
    # Since we don't have any cgo dependencies, force build to be pure
    # golang. This will make multi-os compatibility *MUCH* easier.
    pure = "on",
    visibility = ["//visibility:public"],
    x_defs = {
        "main.Version": "{STABLE_GIT_VERSION}",
    },
)

go_library(
    name = "go_default_library",
    srcs = [
        "fsnotify.go",
        "ibazel.go",
        "lifecycle.go",
        "main.go",
        "source_event_handler.go",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/ibazel",
    visibility = ["//visibility:private"],
    deps = [
        "//bazel:go_default_library",
        "//ibazel/command:go_default_library",
        "//ibazel/live_reload:go_default_library",
        "//ibazel/output_runner:go_default_library",
        "//ibazel/profiler:go_default_library",
        "//ibazel/workspace_finder:go_default_library",
        "//third_party/bazel/master/src/main/protobuf:go_default_library",
        "@com_github_fsnotify_fsnotify//:go_default_library",
    ],
)

go_test(
    name = "go_default_test",
    srcs = [
        "ibazel_test.go",
        "main_test.go",
    ],
    embed = [":go_default_library"],
    importpath = "github.com/bazelbuild/bazel-watcher/ibazel",
    deps = [
        "//bazel:go_default_library",
        "//bazel/testing:go_default_library",
        "//ibazel/command:go_default_library",
        "//ibazel/workspace_finder:go_default_library",
        "//third_party/bazel/master/src/main/protobuf:go_default_library",
        "@com_github_fsnotify_fsnotify//:go_default_library",
        "@com_github_golang_protobuf//proto:go_default_library",
    ],
)
