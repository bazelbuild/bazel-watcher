load("//runfiles_server/rules:serve.bzl", "serve")
load(":integration_test.bzl", "integration_test")
load(":expected_failure_test.bzl", "expected_failure_test")

# First, define a server that serves from the source tree directly.

serve(
    name = "simple",
    data = [
        "hello.txt",
    ],
)

sh_binary(
    name = "simple_test_binary",
    srcs = ["simple_test.sh"],
)

# This test ensures that hello.txt is the same whether read directly from runfiles or indirectly via
# the server.
integration_test(
    name = "simple_test",
    system_under_test = ":simple",
    test_binary = ":simple_test_binary",
)

# Now define a server that serves via a filegroup. It should have the same behavior as :simple.

filegroup(
    name = "filegroup",
    srcs = [
        "hello.txt",
    ],
)

serve(
    name = "simple_filegroup",
    data = [
        ":filegroup",
    ],
)

integration_test(
    name = "simple_filegroup_test",
    system_under_test = ":simple_filegroup",
    test_binary = ":simple_test_binary",
)

# Now define a server that serves a generated output. serve() should work for these too.

genrule(
    name = "hello_generated",
    srcs = [
        "hello.txt",
    ],
    outs = [
        "hello_generated.txt",
    ],
    cmd = "cp $< $@",
)

serve(
    name = "simple_genrule",
    data = [
        ":hello_generated",
    ],
)

sh_binary(
    name = "genrule_test",
    srcs = [
        "genrule_test.sh",
    ],
    data = [
        ":hello_generated",
    ],
)

integration_test(
    name = "simple_genrule_test",
    system_under_test = ":simple_genrule",
    test_binary = ":genrule_test",
)

# Directory traversal tests.
# bazel running a serve() target and requesting a file outside of the target's runfiles should
# return a 404. This is a little tricky to test, because bazel test is sandboxed by default while
# bazel run is not.

# First, create a fake_binary.runfiles directory in bazel-bin. path_traversal_test.sh will attempt
# to read this directory.
# TODO: bazel doesn't know that path_traversal_test has a "dependency" on fake_binary, which is kind
# of the point. But it could lead to spurious passes if this target isn't built before
# path_traversal_test is run.
sh_binary(
    name = "fake_binary",
    srcs = [
        "simple_test.sh",
    ],
    data = [
        "bad.txt",
    ],
)

sh_binary(
    name = "path_traversal_test_binary",
    srcs = [
        "path_traversal_test.sh",
    ],
)

integration_test(
    name = "path_traversal_test_sandboxed",
    system_under_test = ":simple",
    test_binary = ":path_traversal_test_binary",
)

# Even if the devserver allowed path traversal, bazel test :path_traversal_test_sandboxed would pass
# due to sandboxing. Explicitly rerun the test unsandboxed. (Since the serve() rule is designed to
# be used with bazel run, this test is documentation for its behavior with respect to path
# traversal.)
integration_test(
    name = "path_traversal_test_local",
    local = True,
    system_under_test = ":simple",
    test_binary = ":path_traversal_test_binary",
)

# Tests for the integration_test rule itself. An integration_test whose test_binary is fail.sh
# should always fail.

sh_binary(
    name = "fail",
    srcs = [
        "fail.sh",
    ],
)

integration_test(
    name = "should_fail",
    system_under_test = ":simple", # doesn't really matter what it is
    test_binary = ":fail",
    tags = [
        "manual", # don't run this as part of bazel test ... (it's expected to fail)
    ],
)

expected_failure_test(
    name = "expected_failure_test",
    test_binary = ":should_fail",
)