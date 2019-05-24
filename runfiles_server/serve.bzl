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

SERVE_ATTRS = {
    "_server": attr.label(
        default = "//runfiles_server",
        executable = True,
        cfg = "host",
    ),
}

def _serve_impl(ctx):
    return [
        DefaultInfo(
            runfiles = serve_this(ctx, index = ctx.file.index),
        ),
    ]

_serve = rule(
    implementation = _serve_impl,
    doc = """Serves files in an ibazel-aware local development server.
See the serve() macro (which wraps this rule) for documentation.""",
    executable = True,
    attrs = dict(SERVE_ATTRS, **{
        "data": attr.label_list(
            allow_files = True,
        ),
        "index": attr.label(
            allow_single_file = True,
        ),
    }),
)

def serve(name, data=[], index=None, **kwargs):
    """Serves arbitrary files in an ibazel-aware local development server.

    Use this macro when you have a rule whose implementation you do not control that you would like
    to run in a local development server. (If you do control the rule's implementation, use
    the `serve_this` helper function.)

    bazel running a `serve()` target starts a web server on an open port serving the given
    `data` files. If `index` is specified, it also opens the system's default web browser to that
    file.

    When ibazel running a `serve()` target, any changes to the `data` or `index` files are
    immediately displayed in the browser.

    Args:
        name: A unique name for this rule.
        data: Files to serve.
        index: If given, bazel running a serve() target will open a browser pointing to this file.
    """
    _serve(
        name = name,
        data = data,
        index = index,
        # Set the magic flags enabling livereload support.
        # TODO: if ibazel used a different signaling mechanism (a well-known target?),
        # macro wrapping wouldn't be necessary.
        tags = [
            "ibazel_live_reload",
            "ibazel_notify_changes",
        ],
        **kwargs
    )

def serve_this(ctx, index = None, other_files = None):
    """Allows rules to boot an ibazel-aware web server on bazel run.

    This is neither a rule implementation function nor a macro. It is a helper function to be called
    from other rule implementations. Rules producing outputs that can be usefully previewed in a
    browser can call this function to set up all the serving logic. Requirements:

    - The rule must declare `executable = True`, but cannot already produce an executable. (This
       function writes to ctx.outputs.executable.)
    - The rule must mix `SERVE_ATTRS` into its own attributes: `attrs = dict(SERVE_ATTRS, **{...})`.

    Args:
        ctx: The Starlark context object.
        index: The file to open in the system's default browser when the rule is bazel run.
            If not given, bazel run will launch a web server but not the browser.
        other_files: Other files to serve.

    Returns:
        A runfiles object. The calling rule must propagate this as part of its own runfiles.
    """

    # Write a script to invoke the server at bazel run time.
    ctx.actions.write(
        # $@ propagates flags passed to this script to the underlying server binary. Use cases:
        # bazel run <some_target> -- --port 9999 # hard-code port
        # bazel run <some_target> -- --nobrowser # explicitly disable the browser
        content = """#!/bin/sh
%s %s "$@" """ % (
            ctx.executable._server.short_path,
            ("--index " + index.short_path) if index else "",
        ),
        output = ctx.outputs.executable,
    )

    return ctx.runfiles(
        collect_default = True,
        files = [index] if index else [],
        transitive_files = depset(
            transitive = [ctx.attr._server[DefaultInfo].default_runfiles.files] +
                         ([other_files] if other_files else []),
        ),
    )