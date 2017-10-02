# Mock Bazel

This is a special testing implementation of the Bazel CLI API. This is not a go
implementation of the Bazel command line tool. You will only be able to use
this tool to test software that interacts with Bazel over STDOUT/STDIN, signals,
and flag controls.

