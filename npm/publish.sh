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

# Clean the repo to make sure nothing strange is happening.
bazel clean --expunge

echo "Building package for NPM..."
bazel build --config=release "//npm:npm"
echo "Build successful."

echo -n "Publishing ${STAGING} to NPM as "
grep "version" < "$(bazel info bazel-genfiles)/npm/package.json"

# Time to publish...
npm publish "$(bazel info bazel-bin)/npm/npm.tar"
