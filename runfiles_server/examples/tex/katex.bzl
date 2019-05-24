# Copyright 2019 The Bazel Authors. All rights reserved.
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

load("//runfiles_server:serve.bzl", "SERVE_ATTRS", "serve_this")

def _katex_impl(ctx):
    tmp = ctx.actions.declare_file(ctx.attr.name + ".tmp")
    ctx.actions.run_shell(
        inputs = [ctx.file.src],
        outputs = [tmp],
        tools = [ctx.executable._katex],
        command = "< %s %s -d --fleqn > %s" % (
            ctx.file.src.path,
            ctx.executable._katex.path,
            tmp.path,
        ),
        progress_message = "tex -> html %s" % ctx.label,
    )

    ctx.actions.run_shell(
        inputs = [ctx.file._preamble, tmp],
        outputs = [ctx.outputs.html],
        arguments = [ctx.file._preamble.path, tmp.path, ctx.outputs.html.path],
        command = "cat $1 $2 > $3",
    )

    return [
        DefaultInfo(
            runfiles = serve_this(ctx, index = ctx.outputs.html, other_files = ctx.attr._katex_files[DefaultInfo].files),
        ),
    ]

_katex = rule(
    implementation = _katex_impl,
    doc = "runs the katex cli to render .tex files into html at build time",
    executable = True,
    attrs = dict(SERVE_ATTRS, **{
        "src": attr.label(
            mandatory = True,
            allow_single_file = True,
        ),
        "server": attr.label(
            default = "//runfiles_server",
            executable = True,
            cfg = "host",
        ),
        "_katex": attr.label(
            default = "@npm//katex/bin:katex",
            cfg = "host",
            executable = True,
        ),
        "_katex_files": attr.label(
            default = "@npm//katex",
        ),
        "_preamble": attr.label(
            default = ":preamble.html",
            allow_single_file = True,
        ),
    }),
    outputs = {
        "html": "%{name}.html",
    },
)

def katex(name, src):
    """Macro wrapper for _katex, to propagate tags needed for livereload."""
    _katex(
        name = name,
        src = src,
        # Set the magic flags enabling livereload support.
        # TODO: if ibazel used a different signaling mechanism (a well-known target?),
        # macro wrapping wouldn't be necessary.
        tags = [
            "ibazel_live_reload",
            "ibazel_notify_changes",
        ],
    )
