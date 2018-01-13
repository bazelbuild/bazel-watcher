#! /bin/bash

if [ $# -ne "1" ]; then
  cat <<EOF
This script prepares a release for NPM and codifies all the steps required
for tagging a binary.

Usage:

./release.sh tag

Example:

./release.sh v1.0.0

That should tag at version 
EOF
  exit 1
fi

VERSION="$1"; shift

git tag "${VERSION}"

if ./npm/publish.sh; then
  # Success! Publish the tag to GitHub
  git push upstream "${VERSION}"
else
  # Clean up in the event of failure.
  git tag -d "${VERSION}"
fi

