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

# TODO: this just asserts that the call succeeded. More assertions.
curl http://localhost:"$PORT"/hello
