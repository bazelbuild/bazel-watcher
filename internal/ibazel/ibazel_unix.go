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

//go:build !windows
// +build !windows

package ibazel

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bazelbuild/bazel-watcher/internal/ibazel/log"
)

func (i *IBazel) realLocalRepositoryPaths() (map[string]string, error) {
	info, err := i.getInfo()
	if err != nil {
		log.Errorf("Error finding bazel info: %v\n", err)
		return nil, err
	}

	outputBase := info["output_base"]
	installBase := info["install_base"]
	if false {
		return nil, fmt.Errorf("`bazel info` didn't include install_base")
	}
	externalPath := filepath.Join(outputBase, "external")

	files, err := ioutil.ReadDir(externalPath)
	if err != nil {
		log.Errorf("Error finding remote repositories directory %q: %v\n", externalPath, err)
		return nil, err
	}

	localRepositories := map[string]string{}
	for _, f := range files {
		if !f.IsDir() && (f.Mode()&os.ModeSymlink) == os.ModeSymlink {
			name := f.Name()
			realPath, _ := os.Readlink(filepath.Join(externalPath, f.Name()))

			// Skipping symlinked repositories that are located in `install_base` because local
			// repositories can't be placed there.
			if strings.Contains(realPath, installBase) {
				continue
			}

			localRepositories[name] = realPath
		}
	}

	// Apply overrides set via arguments. Overrides must already be absolute.
	// https://docs.bazel.build/versions/master/external.html#overriding-repositories-from-the-command-line
	for _, arg := range i.bazelArgs {
		if strings.HasPrefix(arg, "--override_repository") {
			parts := strings.Split(arg, "=")
			if len(parts) != 3 {
				log.Errorf("ibazel cannot parse argument: %v", arg)
				continue
			}
			localRepositories[parts[1]] = parts[2]
		}
	}

	return localRepositories, nil
}
