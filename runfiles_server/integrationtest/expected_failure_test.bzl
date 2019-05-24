load("@bazel_skylib//lib:shell.bzl", "shell")

def _expected_failure_test(ctx):
    ctx.actions.write(
        content = """#!/bin/sh
if %s; then
    echo "expected %s to have a nonzero exit code, but it exited with $?"
    exit 1
fi""" % (ctx.executable.test_binary.short_path, ctx.executable.test_binary.short_path),
        output = ctx.outputs.executable,
    )
    return [
        DefaultInfo(
            runfiles = ctx.runfiles(
                # Note: expected_failure_tests trivially pass without any runfiles, since executing
                # a nonexistent file returns a nonzero exit code, which the script interprets as
                # success. Consider distinguishing status 127 (file not found).
                transitive_files = depset(
                    transitive = [ctx.attr.test_binary[DefaultInfo].default_runfiles.files],
                ),
            ),
        ),
    ]

expected_failure_test = rule(
    implementation = _expected_failure_test,
    doc = """defines a test that passes if and only if its underlying test_binary fails.
There are ways of testing for expected failures inside unit tests (e.g. junit's ExpectedException).
This rule is useful for testing higher-level integration tests.""",
    test = True,
    attrs = {
        "test_binary": attr.label(
            executable = True,
            cfg = "target",
            mandatory = True,
            doc = """test binary to run. the expected_failure_test will succeed (exit status 0) iff
the test_binary fails (nonzero exit status).""",
        ),
    },
)
