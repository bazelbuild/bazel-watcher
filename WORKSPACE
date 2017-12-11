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

git_repository(
    name = "com_github_bazelbuild_bazel_integration_testing",
    commit = "55a6a70dbcc2cc7699ee715746fb1452788f8d3c",
    remote = "https://github.com/bazelbuild/bazel-integration-testing",
)

load("@com_github_bazelbuild_bazel_integration_testing//tools:repositories.bzl", "bazel_binaries")

bazel_binaries()

rules_go_commit = "dea7cd17fe34744e28e6926feb9efb27ae665b18"

git_repository(
    name = "io_bazel_rules_go",
    commit = rules_go_commit,
    remote = "https://github.com/bazelbuild/rules_go",
)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains", "go_repository")

go_rules_dependencies()

go_register_toolchains()

go_repository(
    name = "com_github_fsnotify_fsnotify",
    commit = "7d7316ed6e1ed2de075aab8dfc76de5d158d66e1",
    importpath = "github.com/fsnotify/fsnotify",
)

go_repository(
    name = "com_github_jaschaephraim_lrserver",
    commit = "50d19f603f71e0a914f23eea33124ba9717e7873",
    importpath = "github.com/jaschaephraim/lrserver",
)

go_repository(
    name = "com_github_gorilla_websocket",
    commit = "7ca4275b84a9d500f68971c8c4a97f0ec18eb889",
    importpath = "github.com/gorilla/websocket",
)

go_repository(
    name = "org_golang_x_sys",
    commit = "99f16d856c9836c42d24e7ab64ea72916925fa97",
    importpath = "golang.org/x/sys",
)

go_repository(
    name = "com_github_bazelbuild_rules_go",
    commit = rules_go_commit,
    importpath = "github.com/bazelbuild/rules_go",
)

go_repository(
    name = "com_github_golang_protobuf",
    commit = "1e59b77b52bf8e4b449a57e6f79f21226d571845",
    importpath = "github.com/golang/protobuf",
)

go_repository(
    name = "com_github_gorilla_websocket",
    commit = "c55883f97322b4bcbf48f734e23d6ab3af1ea488",
    importpath = "github.com/gorilla/websocket",
)
