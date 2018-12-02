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

  DESTINATION="${STAGING}/ibazel_${GOOS}_${GOARCH}"
  if [[ "${GOOS}" == "windows" ]]; then
    DESTINATION="${DESTINATION}.exe"
  fi
  bazel build \
    --config=release \
    "--experimental_platforms=@io_bazel_rules_go//go/toolchain:${GOOS}_${GOARCH}" \
    "//ibazel:ibazel"
  SOURCE="$(bazel info bazel-bin)/ibazel/${GOOS}_${GOARCH}_pure_stripped/ibazel"
  cp "${SOURCE}" "${DESTINATION}"

  # Sometimes bazel likes to change the ouput directory for binaries
  # depending on command line flags (platforms for example). In order to
  # make this an easy to detect error, force remove the binary that was
  # generated for this platform so that if future bazel build runs for a
  # different architecture write to a different folder the expected
  # directory will not exist.
  rm -f "${SOURCE}"
}

# Now compiler ibazel for every platform/arch that is supported.
compile "linux"   "amd64"
compile "darwin"  "amd64"
# Windows isn't compatable due to the os.Setpgid call.
#compile "windows" "amd64"
echo "Build successful."

ghr -t "${CHANGELOG_GITHUB_TOKEN}" "${TAG}" "${STAGING}"
