#!/usr/bin/env bash

set -eo pipefail

PORT=

# Parse the --backend_port flag.
while :; do
  case $1 in
    --backend_port) PORT=$2; shift ;;
    *) break ;;
  esac
  shift
done

# Attempt to read out of another target's runfiles dir. This should fail.
# -f propagates the HTTP status code (404) as a nonzero exit status
# --path-as-is prevents curl from rewriting the URL
# TODO: this binary exits with zero status under bazel test but nonzero status under bazel run
# due to sandboxing. Make sure this test is run without sandboxing too.
! curl -sf --path-as-is http://localhost:"$PORT"/../../fake_binary.runfiles/b/brs/test/bad.txt