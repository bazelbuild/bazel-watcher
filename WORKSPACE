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

load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

git_repository(
    name = "com_github_bazelbuild_bazel_integration_testing",
    commit = "57db374bec74b684e85c5aecde0bf3a26ae27aa7",
    remote = "https://github.com/bazelbuild/bazel-integration-testing",
)

load("@com_github_bazelbuild_bazel_integration_testing//tools:repositories.bzl", "bazel_binaries")

bazel_binaries()

http_archive(
    name = "bazel_skylib",
    sha256 = "b5f6abe419da897b7901f90cbab08af958b97a8f3575b0d3dd062ac7ce78541f",
    strip_prefix = "bazel-skylib-0.5.0",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/archive/0.5.0.tar.gz",
        "https://github.com/bazelbuild/bazel-skylib/archive/0.5.0.tar.gz",
    ],
)

# NOTE: URLs are mirrored by an asynchronous review process. They must
#       be greppable for that to happen. It's OK to submit broken mirror
#       URLs, so long as they're correctly formatted. Bazel's downloader
#       has fast failover.

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "7be7dc01f1e0afdba6c8eb2b43d2fa01c743be1b9273ab1eaf6c233df078d705",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/0.16.5/rules_go-0.16.5.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/0.16.5/rules_go-0.16.5.tar.gz",
    ],
)

http_archive(
    name = "bazel_gazelle",
    sha256 = "6e875ab4b6bf64a38c352887760f21203ab054676d9c1b274963907e0768740d",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/0.15.0/bazel-gazelle-0.15.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/0.15.0/bazel-gazelle-0.15.0.tar.gz",
    ],
)

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

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
    tag = "v1.4.0",
)

go_repository(
    name = "org_golang_x_sys",
    commit = "a5c9d58dba9a56f97aaa86f55e638b718c5a6c42",
    importpath = "golang.org/x/sys",
)

# NOTE: this must match rules_go version located at
# https://github.com/bazelbuild/rules_go/blob/master/go/private/repositories.bzl
go_repository(
    name = "com_github_golang_protobuf",
    importpath = "github.com/golang/protobuf",
    tag = "v1.2.0",
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
    tag = "0.16.2",
)
