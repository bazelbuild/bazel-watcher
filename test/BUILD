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
