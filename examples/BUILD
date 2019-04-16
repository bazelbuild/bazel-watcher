load("//brs/rules:serve.bzl", "serve")
load("@bazel_pandoc//:pandoc.bzl", "pandoc")

serve(
    name = "simple",
    index = "simple.html",
)

# Markdown livereload example. pandoc() is a third-party rule that we cannot change to make
# ibazel-aware. So declare a serve() target that consumes the pandoc() target and is ibazel-aware.

pandoc(
    name = "run_pandoc",
    src = "simple.md",
    from_format = "markdown",
    to_format = "html",
)

serve(
    name = "markdown",
    index = ":run_pandoc",
)
