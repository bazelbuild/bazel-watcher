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
        "@com_google_guava_guava//jar",
        "@flogger//api",
        "@javax_activation//jar",
        "@jcommander//jar",
    ],
)
