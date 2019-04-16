def _serve_impl(ctx):
    # Write a script to invoke the server at bazel run time.
    ctx.actions.write(
        # The $@ propagates flags passed to this executable (ctx.outputs.executable) to the
        # underlying one (ctx.executable._server). This allows the integration test runner to invoke
        # this executable with a --port flag.
        content = '%s %s "$@"' % (
            ctx.executable._server.short_path,
            "--index " + ctx.file.index.short_path if ctx.attr.index else "",
        ),
        output = ctx.outputs.executable,
    )

    transitive_runfiles = [ctx.attr._server[DefaultInfo].default_runfiles.files]
    if ctx.attr.index:
        transitive_runfiles.append(ctx.attr.index[DefaultInfo].files)

    return [
        # Make the data files available at runtime.
        DefaultInfo(
            runfiles = ctx.runfiles(
                collect_default = True,
                transitive_files = depset(
                    transitive = transitive_runfiles,
                ),
            ),
        ),
    ]

_serve = rule(
    implementation = _serve_impl,
    executable = True,
    attrs = {
        "data": attr.label_list(
            allow_files = True,
        ),
        "index": attr.label(
            allow_single_file = True,
        ),
        "_server": attr.label(
            default = "//brs",
            executable = True,
            cfg = "host",
        ),
    },
)

def serve(name, data = None, index = None):
    """Macro wrapper for _serve, to propagate tags needed for livereload."""
    _serve(
        name = name,
        data = data or [],
        index = index,
        tags = [
            "ibazel_live_reload",
            "ibazel_notify_changes",
        ],
    )
