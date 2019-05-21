load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_binary(
    name = "main",
    srcs = [
        "RunfilesServer.go",
    ],
    visibility = [
        "//visibility:public",
    ],
)

java_library(
    name = "EphemeralPort",
    srcs = [
        "EphemeralPort.java",
    ],
    visibility = [
        "//visibility:public",
    ],
)

java_binary(
    name = "brs",
    srcs = [
        "RunfilesServer.java",
    ],
    data = [
        "@apache_mime_types//file",
    ],
    main_class = "brs.RunfilesServer",
    visibility = [
        "//visibility:public",
    ],
    runtime_deps = [
        "@flogger//google:flogger",
    ],
    deps = [
        ":EphemeralPort",
        "@com_google_code_findbugs_jsr305//jar",
        "@com_google_guava_guava//jar",
        "@flogger//api",
        "@javax_activation//jar",
        "@jcommander//jar",
    ],
)

java_binary(
    name = "IntegrationTestRunner",
    testonly = True,
    srcs = [
        "IntegrationTestRunner.java",
    ],
    main_class = "brs.IntegrationTestRunner",
    visibility = [
        "//visibility:public",
    ],
    runtime_deps = [
        "@flogger//google:flogger",
    ],
    deps = [
        ":EphemeralPort",
        "@com_google_guava_guava//jar",
        "@flogger//api",
        "@jcommander//jar",
    ],
)
