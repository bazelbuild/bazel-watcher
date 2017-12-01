#! /usr/bin/env bash

echo -n "STABLE_GIT_VERSION "
if git diff-index --quiet HEAD -- > /dev/null 2>&1; then
  git describe --tags --abbrev=0
else
  printf "%s-dirty\n" "$(git describe --tags --abbrev=0)"
fi

