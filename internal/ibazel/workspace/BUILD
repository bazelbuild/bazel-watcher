# Copyright 2018 The Bazel Authors. All rights reserved.
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
    name = "workspace",
    srcs = ["workspace.go"],
    importpath = "github.com/bazelbuild/bazel-watcher/internal/ibazel/workspace",
    visibility = ["//:__subpackages__"],
    deps = ["//internal/ibazel/log"],
)

go_test(
    name = "workspace_test",
    srcs = ["workspace_test.go"],
    embed = [":workspace"],
    deps = [
        "//internal/ibazel/log",
        "@com_github_google_go_cmp//cmp",
    ],
)
