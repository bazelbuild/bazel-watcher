# Bazel runfiles server

This repo provides a web server for local development that serves directly out of a Bazel target's
runfiles with ~zero extra semantics (no bundling, minification, etc.). It also provides Starlark
definitions that wrap the server, allowing any rule to integrate with ibazel livereload with minimal
changes. This enables a powerful local development experience, with sub-second preview latency in
some cases.

## Background

[ibazel](https://github.com/bazelbuild/bazel-watcher) is a valuable development tool. When an executable target is run with ibazel and the target's
inputs change, ibazel kills whatever the running target was doing and restarts it with the changed
inputs. For rules that do something simple like print to a terminal, this is an ideal development
experience, giving immediate feedback.

Other rules produce outputs that can be usefully previewed in a browser (HTML/CSS/JS, but also
Markdown, SVG, TeX, ...). For these, the ideal ibazel development experience is to preview them in
a browser, reloading the page automatically whenever the target's inputs change. This is currently
possible thanks to ibazel's livereload support, but integrating it into custom rules is complex.
Adding logic to a rule to boot a web server and launch a browser distracts from the rule's main
business logic.

My vision is this: **any rule that produces outputs that are valuable to show in a browser should be
able to be invoked with ibazel run and have livereload just work**.

## Usage

The repo provides hooks for both rule authors and rule users. In either case, you first need to add
the following to your WORKSPACE to set up the dependencies:

```py
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

git_repository(
  name = "com_github_ubehebe_bazel_runfiles_server",
  remote = "https://github.com/Ubehebe/bazel-runfiles-server.git",
  commit = ... # TODO: add a valid commit here
)

load("@com_github_ubehebe_bazel_runfiles_server//:repositories.bzl", "bazel_runfiles_server_repositories")

bazel_runfiles_server_repositories()

load("@com_github_ubehebe_bazel_runfiles_server//:setup.bzl", "setup")

setup()
```

### For rule authors

[serve_this()](rules/serve.bzl#L45) is a Starlark helper function that you can drop into your rule
implementation to make it livereload-aware. It returns a runfiles object that you propagate as the
runfiles of your rule:

```py
load("@com_github_ubehebe_bazel_runfiles_server//rules:serve.bzl", "serve_this")

def foo_impl(ctx):
  index = ... # generate output that the browser will open by default when this target is bazel run
  other_files = ... # generate other outputs to serve (CSS, JS, etc.)
  return [
    DefaultInfo(runfiles=serve_this(ctx, index, other_files))
  ]
```

For example, [examples/tex/katex.bzl](examples/tex/katex.bzl) is a rule that uses the
[KaTeX](https://katex.org) command-line tool to render TeX into HTML at build time. The use of
serve_this allows `katex()` targets to be directly runnable with bazel and ibazel, and to support
livereload when run with the latter.

### For rule users

If you are using a rule that you would like to be livereload-aware, but you do not control its
implementation, use the [serve()](rules/serve.bzl#L30) macro in your BUILD file, and pass the
underlying target in as a data or index dependency:

```py
load("@com_github_ubehebe_bazel_runfiles_server//rules:serve.bzl", "serve")

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

For example, [//examples:run_pandoc](examples/BUILD#L19) uses the third-party
[Pandoc rules](https://github.com/ProdriveTechnologies/bazel-pandoc) to render Markdown into HTML at
build time. Since I don't control the pandoc() rule implementation, I define a serve() target
([//examples:markdown](examples/BUILD#L26)) that consumes the pandoc() rule.

serve() can also be used to serve files directly from the source tree, from filegroups, from genrule
outputs. etc.