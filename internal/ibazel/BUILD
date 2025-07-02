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

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "ibazel",
    srcs = [
        "ibazel.go",
        "ibazel_unix.go",
        "ibazel_windows.go",
        "lifecycle.go",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/internal/ibazel",
    visibility = ["//:__subpackages__"],
    deps = [
        "//internal/bazel",
        "//internal/ibazel/command",
        "//internal/ibazel/fswatcher",
        "//internal/ibazel/fswatcher/common",
        "//internal/ibazel/lifecycle_hooks",
        "//internal/ibazel/live_reload",
        "//internal/ibazel/log",
        "//internal/ibazel/output_runner",
        "//internal/ibazel/profiler",
        "//internal/ibazel/workspace",
        "//third_party/bazel/master/src/main/protobuf/blaze_query",
    ],
)

go_test(
    name = "ibazel_test",
    srcs = ["ibazel_test.go"],
    embed = [":ibazel"],
    importpath = "github.com/bazelbuild/bazel-watcher/ibazel",
    deps = [
        "//internal/bazel",
        "//internal/bazel/testing",
        "//internal/ibazel/command",
        "//internal/ibazel/fswatcher/common",
        "//internal/ibazel/log",
        "//internal/ibazel/workspace",
        "//third_party/bazel/master/src/main/protobuf/analysis",
        "//third_party/bazel/master/src/main/protobuf/blaze_query",
        "@com_github_golang_protobuf//proto",
    ],
)
