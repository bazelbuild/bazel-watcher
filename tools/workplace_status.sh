#! /usr/bin/env bash

echo -n "STABLE_GIT_VERSION "

# A hack to make the `git describe` work correctly in Travis CI in forks that
# are not the original bazelbuild repo. This assumes artifacts are actually
# built using the ci.bazelbuild jobs.
if [[ $TRAVIS == "true" ]]; then
  printf "%s-dirty\n" "$(git rev-parse HEAD)"
  exit
fi

if git diff-index --quiet HEAD -- > /dev/null 2>&1; then
  git describe --tags --abbrev=0
else
  printf "%s-dirty\n" "$(git describe --tags --abbrev=0)"
fi

