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

load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

genrule(
    name = "package",
    srcs = [
        "//:CONTRIBUTORS",
        ":npm",
    ],
    outs = ["package.json"],
    cmd = "$(location :package_sh) ./$(location :npm) $(location //:CONTRIBUTORS) > $@",
    stamp = 1,
    tags = [
        "manual",
        "nowindows",
    ],
    tools = [
        ":package_sh",
    ],
    visibility = ["//visibility:private"],
)

sh_binary(
    name = "package_sh",
    srcs = ["package.sh"],
    visibility = ["//visibility:private"],
)

go_library(
    name = "npm_lib",
    srcs = ["generate_package.go"],
    importpath = "github.com/bazelbuild/bazel-watcher/release/npm",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "npm",
    embed = [":npm_lib"],
    importpath = "github.com/bazelbuild/bazel-watcher/npm",
    tags = [
        "manual",
    ],
    visibility = ["//visibility:private"],
    x_defs = {
        "main.Version": "{STABLE_GIT_VERSION}",
    },
)
