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

readonly TAG="$1"; shift

if !git rev-parse "${TAG}" >/dev/null 2>&1; then
  echo "The provided tag doesn't exist. First create, then push the tag"
  exit 1
fi

# Make a temporary directory to stage the release in.
readonly STAGING="$(mktemp -d)"
echo "Staging into ${STAGING}"

compile() {
  export GOOS=$1; shift
  export GOARCH=$1; shift
  export CGO=$1; shift

  EXTENSION=""
  if [[ "${GOOS}" == "windows" ]]; then
    EXTENSION=".exe"
  fi
  PURE="pure_"
  TOOLCHAIN="${GOOS}_${GOARCH}"
  if [[ "${CGO}" == "cgo" ]]; then
    PURE=""
    TOOLCHAIN="${TOOLCHAIN}_cgo"
  fi
  DESTINATION="${STAGING}/ibazel_${GOOS}_${GOARCH}${EXTENSION}"
  bazel build \
    --config=release \
    "--platforms=@io_bazel_rules_go//go/toolchain:${TOOLCHAIN}" \
    "//ibazel:ibazel"
  SOURCE="$(bazel info bazel-bin)/ibazel/${GOOS}_${GOARCH}_${PURE}stripped/ibazel${EXTENSION}"
  cp "${SOURCE}" "${DESTINATION}"

  # Sometimes bazel likes to change the ouput directory for binaries
  # depending on command line flags (platforms for example). In order to
  # make this an easy to detect error, force remove the binary that was
  # generated for this platform so that if future bazel build runs for a
  # different architecture write to a different folder the expected
  # directory will not exist.
  rm -f "${SOURCE}"

  # Clean out the old environment variables so that you don't build GHR for
  # windows on accident.
  unset GOOS
  unset GOARCH
}

# Now compiler ibazel for every platform/arch that is supported.
compile "linux"   "amd64"  ""
compile "darwin"  "amd64"  "cgo"
compile "windows" "amd64"  ""

echo "Build successful."

readonly GHR_BINARY="$(mktemp /tmp/ghr.XXXXXX)"
go get -u github.com/tcnksm/ghr
go build -o "${GHR_BINARY}" github.com/tcnksm/ghr
chmod +x "${GHR_BINARY}"
echo "Publishing ${STAGING} to GitHub as ${TAG}"
"${GHR_BINARY}" \
  -t "${CHANGELOG_GITHUB_TOKEN}" \
  -u "bazelbuild" \
  -r "bazel-watcher" \
  "${TAG}" \
  "${STAGING}"
find "${STAGING}"
