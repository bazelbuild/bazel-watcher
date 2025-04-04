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
    name = "fsevents",
    srcs = ["fsevents.go"],
    importpath = "github.com/bazelbuild/bazel-watcher/internal/ibazel/fswatcher/fsevents",
    visibility = ["//:__subpackages__"],
    deps = select({
        "@io_bazel_rules_go//go/platform:darwin": [
            "//internal/ibazel/fswatcher/common",
            "@com_github_fsnotify_fsevents//:fsevents",
        ],
        "@io_bazel_rules_go//go/platform:ios": [
            "//internal/ibazel/fswatcher/common",
            "@com_github_fsnotify_fsevents//:fsevents",
        ],
        "//conditions:default": [],
    }),
)

go_test(
    name = "fsevents_test",
    srcs = ["fsevents_test.go"],
    embed = [":fsevents"],
    importpath = "github.com/bazelbuild/bazel-watcher/internal/ibazel/fswatcher/fsevents",
    deps = select({
        "@io_bazel_rules_go//go/platform:darwin": [
            "@com_github_google_go_cmp//cmp",
        ],
        "@io_bazel_rules_go//go/platform:ios": [
            "@com_github_google_go_cmp//cmp",
        ],
        "//conditions:default": [],
    }),
)
