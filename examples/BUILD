load("//brs/examples/tex:katex.bzl", "katex")
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

# TeX livereload example. katex() is a rule defined in this repo. Its implementation calls the
# serve_this() helper function, making all katex() targets directly runnable and ibazel-aware.
# No need to define a separate serve() target.

katex(
    name = "tex",
    src = "simple.tex",
)
