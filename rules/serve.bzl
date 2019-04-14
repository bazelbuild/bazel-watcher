def _serve(ctx):
    # Write a script to invoke the server at bazel run time.
    ctx.actions.write(
        # The $@ propagates flags passed to this executable (ctx.outputs.executable) to the
        # underlying one (ctx.executable._server). This allows the integration test runner to invoke
        # this executable with a --port flag.
        content = '%s "$@"' % ctx.executable._server.short_path,
        output = ctx.outputs.executable,
    )

    return [
        # Make the data files available at runtime.
        DefaultInfo(
            runfiles = ctx.attr._server[DefaultInfo].default_runfiles.merge(
                ctx.runfiles(collect_default = True),
            ),
        ),
    ]

serve = rule(
    implementation = _serve,
    executable = True,
    attrs = {
        "data": attr.label_list(
            allow_files = True,
        ),
        "_server": attr.label(
            default = "//brs",
            executable = True,
            cfg = "host",
        ),
    },
)
