# Bazel watcher

[![Build Status](https://ci.bazel.io/buildStatus/icon?job=Global%2Fbazel-watcher)](https://ci.bazel.io/blue/organizations/jenkins/Global%2Fbazel-watcher/activity/)

Note: This is not an official Google product.

A source file watcher for [Bazel](https://Bazel.build) projects

Ever wanted to save a file and have your tests automatically run? How about
restart your webserver when one of the source files change? Look no further.

Compile the `//ibazel` target inside this repo and copy the source file onto
your `$PATH`.

Then:

```bash
# ibazel build //path/to/my:target
```

Hack hack hack. Save and your target will be rebuilt.

Right now this repo supports `build`, `test`, and `run`.

## Additional notes

### Termination

SIGINT has to be sent twice to kill ibazel: once to kill the subprocess, and
the second time for ibazel itself. Also, ibazel will exit on its own when a
bazel query fails, but it will stay alive when a build, test, or run fails.
We use an exit code of 3 for a signal termination, and 4 for a query failure,
but these codes are subject to change.

### What about the `--watchfs` flag?

Bazel has a flag called `--watchfs` which, according to the bazel command-line
help does:

> If true, Bazel tries to use the operating system's file watch service for
> local changes instead of scanning every file for a change

Unfortunately, this option does not rebuild the project on save like the Bazel
watcher does, but instead queries the file system for a list of files that have
been invalidated since last build and will require reinspection by the Bazel
server.

Copyright 2017 The Bazel Authors. All right reserved.
