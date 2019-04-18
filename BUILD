load("//rules:serve.bzl", "serve")
load("@bazel_pandoc//:pandoc.bzl", "pandoc")

# A nice livereload experience for writing the readme itself.

pandoc(
    name = "run_pandoc",
    src = "README.md",
    from_format = "markdown",
    to_format = "html",
)

serve(
    name = "readme",
    index = ":run_pandoc",
)
