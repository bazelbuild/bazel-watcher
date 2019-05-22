def _integration_test(ctx):
    ctx.actions.write(
        content = "%s --sut_binary %s --test_binary %s" % (
            ctx.executable._test_runner.short_path,
            ctx.executable.system_under_test.short_path,
            ctx.executable.test_binary.short_path,
        ),
        output = ctx.outputs.executable,
    )

    return [
        DefaultInfo(
            runfiles = ctx.runfiles(
                transitive_files = depset(
                    transitive = [
                        ctx.attr.system_under_test[DefaultInfo].default_runfiles.files,
                        ctx.attr.test_binary[DefaultInfo].default_runfiles.files,
                        ctx.attr._test_runner[DefaultInfo].default_runfiles.files,
                    ],
                ),
            ),
        ),
    ]

integration_test = rule(
    implementation = _integration_test,
    test = True,
    attrs = {
        "system_under_test": attr.label(
            mandatory = True,
            executable = True,
            cfg = "target",
        ),
        "test_binary": attr.label(
            executable = True,
            mandatory = True,
            cfg = "target",
        ),
        "_test_runner": attr.label(
            default = "//brs:IntegrationTestRunner",
            executable = True,
            cfg = "target",
        ),
    },
)
