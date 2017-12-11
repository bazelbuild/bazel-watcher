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
	"errors"
	"os"
	"path/filepath"
)

type WorkspaceFinder interface {
    FindWorkspace() (string, error)
}

type MainWorkspaceFinder struct {}

func (m *MainWorkspaceFinder) FindWorkspace() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	volume := filepath.VolumeName(path)

	for {
		// filepath.Dir() includes a trailing separator if we're at the root
		if path == volume + string(filepath.Separator) {
			path = volume
		}

		// Check if we're at the workspace path
		if _, err := os.Stat(filepath.Join(path, "WORKSPACE")); !os.IsNotExist(err) {
			return path, nil
		}

		// If we've reached the root, then we know the cwd isn't within a workspace
		if path == volume {
			return "", errors.New("ibazel was not invoked from within a workspace\n")
		}

		// Move one level up the path
		path = filepath.Dir(path)
	}
}

type FakeWorkspaceFinder struct {}

func (f *FakeWorkspaceFinder) FindWorkspace() (string, error) {
    return "", nil
}
