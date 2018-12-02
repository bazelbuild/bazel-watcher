#!/usr/bin/env bash

readonly WORKSPACE="$(bazel info workspace)"

cd "${WORKSPACE}"

if [ $# -ne "1" ]; then
  cat <<EOF
This script prepares a release and codifies all the steps required for tagging
a binary.

Usage:

./release.sh tag

Example:

./release.sh v1.0.0

That should tag at version
EOF
  exit 1
fi

set -ex

GIT_TAG=$1; shift

cat - <<EOF
Releasing bazel-watcher at tag ${GIT_TAG}.

This will update the changelog, tag a release and release it on NPM.

Press enter to continue
EOF

read

docker run --rm \
  --interactive \
  --tty \
  -v "${PWD}:/usr/local/src/your-app" \
  -e "CHANGELOG_GITHUB_TOKEN=${CHANGELOG_GITHUB_TOKEN}" \
  ferrarimarco/github-changelog-generator:1.14.3 \
      -u bazelbuild \
      -p bazel-watcher \
      --author \
      --compare-link \
      --github-site=https://github.com/bazelbuild/bazel-watcher \
      --unreleased-label "**Next release**" \
      --future-release="${GIT_TAG}"

# Add the newly generated changelog and commit it.
git add CHANGELOG.md
git commit -m "Generating CHANGELOG.md for release ${GIT_TAG}"

# Tag the release.
git tag "${GIT_TAG}"

if ./npm/publish.sh; then
  # Success! Publish the tag to GitHub
  git push git@github.com:bazelbuild/bazel-watcher "${GIT_TAG}"

  # Advance master branch to the tag.
  git push git@github.com:bazelbuild/bazel-watcher "${GIT_TAG}:master"

  # Now the tag exists on Github, we can publish assets to it.
  ./release/github/publish.sh "${GIT_TAG}"
else
  # Clean up in the event of failure.
  git tag -d "${GIT_TAG}"
fi
