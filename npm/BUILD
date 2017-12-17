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

load("//:os.bzl", "OSES")
load("@bazel_tools//tools/build_defs/pkg:pkg.bzl", "pkg_tar")
load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

genrule(
    name = "package",
    srcs = [
        "//:CONTRIBUTORS",
    ],
    outs = ["package.json"],
    cmd = "$(location :package_generator) $(location //:CONTRIBUTORS) > $@",
    stamp = 1,
    tools = [
        ":package_generator",
    ],
)

go_library(
    name = "go_default_library",
    srcs = ["generate_package.go"],
    importpath = "github.com/bazelbuild/bazel-watcher/npm",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "package_generator",
    embed = [":go_default_library"],
    importpath = "github.com/bazelbuild/bazel-watcher/npm",
    visibility = ["//visibility:public"],
    x_defs = {
        "main.Version": "{STABLE_GIT_VERSION}",
    },
)

genrule(
    name = "index_for_oses",
    srcs = ["index.tpl.js"],
    outs = ["index.js"],
    cmd = """
        sed -e 's|__OSES__|%s|' $(location :index.tpl.js) > $(location :index.js)
    """ % ("\", ".join(OSES)),
)

pkg_tar(
    name = "npm_tar",
    srcs = [
        ":index.js",
        ":package.json",
    ],
    mode = "0644",
    package_dir = "/",
    strip_prefix = "/npm",
)

[go_binary(
    name = "ibazel_%s" % GOOS,
    embed = ["//ibazel:go_default_library"],
    goos = GOOS,
    importpath = "github.com/bazelbuild/bazel-watcher/ibazel",
    # Since we don't have any cgo dependencies, force build to be pure
    # golang. This will make multi-os compatability *MUCH* easier.
    pure = "on",
    visibility = ["//visibility:private"],
    x_defs = {
        "main.Version": "{STABLE_GIT_VERSION}",
    },
) for GOOS in OSES]

[pkg_tar(
    name = "%s_tar" % GOOS,
    srcs = [":ibazel_%s" % GOOS],
    mode = "0755",
    package_dir = "/bin/%s_amd64" % GOOS,
    strip_prefix = "/ibazel/%s_amd64_pure_stripped" % GOOS,
) for GOOS in OSES]

pkg_tar(
    name = "npm",
    extension = "tar",
    package_dir = "/",
    deps = [
        ":npm_tar",
    ] + [
        "%s_tar" % GOOS
        for GOOS in OSES
    ],
)
