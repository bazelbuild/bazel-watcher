bazel build @io_bazel//src/tools/skylark/java/com/google/devtools/skylark/skylint:Skylint &&
  find . -type f -name BUILD | 
  xargs $(bazel info bazel-bin)/external/io_bazel/src/tools/skylark/java/com/google/devtools/skylark/skylint/Skylint
