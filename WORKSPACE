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

http_archive(
    name = "com_github_bazelbuild_bazel_integration_testing",
    sha256 = "81a2ad3a8ec5a9d1d91b9aca0b4f1f3a0b094f30c48d582e5226defccd714bb9",
    strip_prefix = "bazel-integration-testing-404010b3763262526d3a0e09073d8a8f22ed3d4b",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-integration-testing/archive/404010b3763262526d3a0e09073d8a8f22ed3d4b.tar.gz",
        "https://github.com/bazelbuild/bazel-integration-testing/archive/404010b3763262526d3a0e09073d8a8f22ed3d4b.tar.gz",
    ],
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
    sha256 = "f5127a8f911468cd0b2d7a141f17253db81177523e4429796e14d429f5444f5f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/0.16.1/rules_go-0.16.1.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/0.16.1/rules_go-0.16.1.tar.gz",
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

# NOTE: this must match rules_go version located at
# https://github.com/bazelbuild/rules_go/blob/master/go/private/repositories.bzl
go_repository(
    name = "com_github_golang_protobuf",
    commit = "b4deda0973fb4c70b50d226b1af49f3da59f5265",
    importpath = "github.com/golang/protobuf",
)

go_repository(
    name = "com_github_gorilla_websocket",
    commit = "c55883f97322b4bcbf48f734e23d6ab3af1ea488",
    importpath = "github.com/gorilla/websocket",
)

# NOTE: this must match rules_go version above, currently set to 0.15.3
go_repository(
    name = "com_github_bazelbuild_rules_go",
    commit = "0f0d007c89dc67a5a34490acafc5195b191f5045",
    importpath = "github.com/bazelbuild/rules_go",
)
