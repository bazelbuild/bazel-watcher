load("//rules:integration_test.bzl", "integration_test")

sh_binary(
    name = "golang_sanity_test_sh",
    srcs = [
        "golang_sanity_test.sh",
    ],
)

integration_test(
    name = "golang_sanity_test",
    system_under_test = "//brs:main",
    test_binary = ":golang_sanity_test_sh",
)