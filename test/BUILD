load("//brs/rules:serve.bzl", "serve")
load("//brs/rules:integration_test.bzl", "integration_test")

serve(
    name = "simple",
    data = [
        "hello.txt",
    ],
)

sh_binary(
    name = "simple_test_binary",
    srcs = [
        "simple_test.sh",
    ],
    data = [
        "hello.txt",
    ],
)

integration_test(
    name = "simple_test",
    system_under_test = ":simple",
    test_binary = ":simple_test_binary",
)

sh_binary(
    name = "path_traversal_test_binary",
    srcs = [
        "path_traversal_test.sh",
    ],
)

# This target creates a fake_binary.runfiles directory in bazel-bin, which path_traversal_test.sh
# attempts to read.
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

integration_test(
    name = "path_traversal_test",
    system_under_test = ":simple",
    test_binary = ":path_traversal_test_binary",
)
