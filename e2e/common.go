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

package e2e

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

func GetPath(p string) string {
	path, err := bazel.Runfile(p)
	if err != nil {
		panic(err)
	}

	path, err = filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return path
}

var ibazelPath string

func init() {
	var err error
	ibazelPath = GetPath(fmt.Sprintf("ibazel/%s_%s_pure_stripped/ibazel", runtime.GOOS, runtime.GOARCH))
	if err != nil {
		panic(err)
	}
}
