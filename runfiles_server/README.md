# Bazel runfiles server

This package provides a web server for local development that serves directly out of a Bazel
target's runfiles with ~zero extra semantics (no bundling, minification, etc.). It also provides
Starlark definitions that wrap the server, allowing any rule to integrate with ibazel livereload
with minimal changes. This enables a powerful local development experience, with sub-second preview
latency in some cases.

## Usage

This package provides hooks for both rule authors and rule users. In either case, add the following
to your WORKSPACE to set up the dependencies:

```py
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

git_repository(
  name = "ibazel",
  remote = "https://github.com/bazelbuild/bazel-watcher.git",
  commit = ... # TODO: add a valid commit here
)

load("@ibazel//:repositories.bzl", "bazel_runfiles_server_repositories")

bazel_runfiles_server_repositories()

load("@ibazel//:setup.bzl", "setup")

setup()
```

### For rule authors

[serve_this()](rules/serve.bzl#L45) is a Starlark helper function that you can drop into your rule
implementation to make it livereload-aware. It returns a runfiles object that you propagate as the
runfiles of your rule:

```py
load("@ibazel//runfiles_server:serve.bzl", "serve_this")

def foo_impl(ctx):
  index = ... # generate output that the browser will open by default when this target is bazel run
  other_files = ... # generate other outputs to serve (CSS, JS, etc.)
  return [
    DefaultInfo(runfiles=serve_this(ctx, index, other_files))
  ]
```

For example, [runfiles_server/examples/tex/katex.bzl](runfiles_server.examples/tex/katex.bzl) is a
rule that uses the [KaTeX](https://katex.org) command-line tool to render TeX into HTML at build
time. The use of serve_this allows `katex()` targets to be directly runnable with bazel and ibazel,
and to support livereload when run with the latter.

### For rule users

If you are using a rule that you would like to be livereload-aware, but you do not control its
implementation, use the [serve()](rules/serve.bzl#L30) macro in your BUILD file, and pass the
underlying target in as a data or index dependency:

```py
load("@ibazel//runfiles_server:serve.bzl", "serve")

foo_library(
  name = "foo",
  ...
)

serve(
  name = "server",
  data = [":foo"],
)
```

The `server` target is runnable with bazel and ibazel, and supports livereload when run with the
latter.

For example, [//runfiles_server/examples:run_pandoc](examples/BUILD#L19) uses the third-party
[Pandoc rules](https://github.com/ProdriveTechnologies/bazel-pandoc) to render Markdown into HTML at
build time. Since we don't control the pandoc() rule implementation, we define a serve()
target ([//runfiles_server/examples:markdown](examples/BUILD#L26)) that consumes the pandoc() rule.

serve() can also be used to serve files directly from the source tree, from filegroups, from genrule
outputs. etc.