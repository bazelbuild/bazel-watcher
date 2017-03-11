# Bazel watcher

Note: This is not an official Google product.

A source file watcher for [Bazel.io](https://Bazel.io) projects

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

Copyright 2017 The Bazel Authors. All right reserved.
