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
    doc = """Brings up a system under test, then runs a test binary against it.

This rule is for testing the serve() rule. It is not generalizable for testing arbitrary servers.

The system_under_test binary will be invoked with a --port flag giving the port
that the system under test should listen on.

The test binary will be invoked with a --backend_port flag giving the port of the system under test.
The test binary should make calls to the system under test using this port.

The exit code of the test_binary determines the overall result of the integration test (zero for
success, nonzero for failure).""",
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
            default = "//runfiles_server:IntegrationTestRunner",
            executable = True,
            cfg = "target",
        ),
    },
)
