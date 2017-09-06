# Copyright 2017 The Bazel Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -o errexit
set -o errtrace
set -o noclobber
set -o pipefail

# cd to the root of the project so that relative paths work.
cd "$(git rev-parse --show-toplevel)"

# Make a temporary directory to stage the release in.
readonly STAGING="$(mktemp -d)"
echo "Staging into ${STAGING}"

# Copy over the base files required for NPM
cp "README.md" "${STAGING}/README.md"
cp "npm/index.js" "${STAGING}/index.js"
cp "npm/package.json" "${STAGING}/package.json"

compile() {
  export GOOS=$1; shift
  export GOARCH=$1; shift

  mkdir -p "${STAGING}/bin/${GOOS}_${GOARCH}/"
  DESTINATION="${STAGING}/bin/${GOOS}_${GOARCH}/ibazel"
  if [[ "${GOOS}" == "windows" ]]; then
    DESTINATION="${DESTINATION}.exe"
  fi
  go build -o "${DESTINATION}" github.com/bazelbuild/bazel-watcher/ibazel
}

# Now compiler ibazel for every platform/arch that is supported.
compile "linux"   "amd64"
compile "darwin"  "amd64"
compile "windows" "amd64"

# Everything is staged now, actually upload the package.
cd "$STAGING" && npm publish
