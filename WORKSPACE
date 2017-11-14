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

rules_go_commit = "0bd97fc6ae48d4124e1a8506bcb397be766f8b83"

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
    commit = "130e6b02ab059e7b717a096f397c5b60111cae74",
    importpath = "github.com/golang/protobuf",
)
