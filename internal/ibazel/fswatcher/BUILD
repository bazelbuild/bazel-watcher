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

load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "fswatcher",
    srcs = [
        "factory_darwin.go",
        "factory_fsnotify.go",
    ],
    importpath = "github.com/bazelbuild/bazel-watcher/internal/ibazel/fswatcher",
    visibility = ["//:__subpackages__"],
    deps = select({
        "@io_bazel_rules_go//go/platform:aix": [
            "//internal/ibazel/fswatcher/common",
            "//internal/ibazel/fswatcher/fsnotify",
        ],
        "@io_bazel_rules_go//go/platform:android": [
            "//internal/ibazel/fswatcher/common",
            "//internal/ibazel/fswatcher/fsnotify",
        ],
        "@io_bazel_rules_go//go/platform:darwin": [
            "//internal/ibazel/fswatcher/common",
            "//internal/ibazel/fswatcher/fsevents",
            "//internal/ibazel/fswatcher/fsnotify",
        ],
        "@io_bazel_rules_go//go/platform:dragonfly": [
            "//internal/ibazel/fswatcher/common",
            "//internal/ibazel/fswatcher/fsnotify",
        ],
        "@io_bazel_rules_go//go/platform:freebsd": [
            "//internal/ibazel/fswatcher/common",
            "//internal/ibazel/fswatcher/fsnotify",
        ],
        "@io_bazel_rules_go//go/platform:illumos": [
            "//internal/ibazel/fswatcher/common",
            "//internal/ibazel/fswatcher/fsnotify",
        ],
        "@io_bazel_rules_go//go/platform:ios": [
            "//internal/ibazel/fswatcher/common",
            "//internal/ibazel/fswatcher/fsevents",
            "//internal/ibazel/fswatcher/fsnotify",
        ],
        "@io_bazel_rules_go//go/platform:js": [
            "//internal/ibazel/fswatcher/common",
            "//internal/ibazel/fswatcher/fsnotify",
        ],
        "@io_bazel_rules_go//go/platform:linux": [
            "//internal/ibazel/fswatcher/common",
            "//internal/ibazel/fswatcher/fsnotify",
        ],
        "@io_bazel_rules_go//go/platform:netbsd": [
            "//internal/ibazel/fswatcher/common",
            "//internal/ibazel/fswatcher/fsnotify",
        ],
        "@io_bazel_rules_go//go/platform:openbsd": [
            "//internal/ibazel/fswatcher/common",
            "//internal/ibazel/fswatcher/fsnotify",
        ],
        "@io_bazel_rules_go//go/platform:osx": [
            "//internal/ibazel/fswatcher/common",
            "//internal/ibazel/fswatcher/fsnotify",
        ],
        "@io_bazel_rules_go//go/platform:plan9": [
            "//internal/ibazel/fswatcher/common",
            "//internal/ibazel/fswatcher/fsnotify",
        ],
        "@io_bazel_rules_go//go/platform:qnx": [
            "//internal/ibazel/fswatcher/common",
            "//internal/ibazel/fswatcher/fsnotify",
        ],
        "@io_bazel_rules_go//go/platform:solaris": [
            "//internal/ibazel/fswatcher/common",
            "//internal/ibazel/fswatcher/fsnotify",
        ],
        "@io_bazel_rules_go//go/platform:windows": [
            "//internal/ibazel/fswatcher/common",
            "//internal/ibazel/fswatcher/fsnotify",
        ],
        "//conditions:default": [],
    }),
)
