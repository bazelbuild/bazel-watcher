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

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "process_group",
    srcs = [
        "process_group.go",
        "process_group_unix.go",
        "process_group_windows.go",
        "syscalls_windows.go",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/internal/ibazel/process_group",
    visibility = ["//:__subpackages__"],
)
