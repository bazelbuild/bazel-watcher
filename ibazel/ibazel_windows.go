// Copyright 2017 The Bazel Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"github.com/bazelbuild/bazel-watcher/ibazel/log"
)

var alreadyNotifiedOfLocalRepositories bool

func (i *IBazel) realLocalRepositoryPaths() (map[string]string, error) {
	if !alreadyNotifiedOfLocalRepositories {
		// Put the entire implementation of this method in here
		log.Banner(
			"iBazel does not support watching local_repository or --override_repository",
			"If this is a feature you'd like to add support for, please visit",
			"https://github.com/bazelbuild/bazel-watcher/issues/274",
		)
	}
	alreadyNotifiedOfLocalRepositories = true
	return map[string]string{}, nil
}
