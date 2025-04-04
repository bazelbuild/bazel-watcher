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

package workspace

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bazelbuild/bazel-watcher/internal/ibazel/log"
)

type Workspace interface {
	FindWorkspace() (string, error)
	ExecuteCommand(command string, args []string)
}

type MainWorkspace struct{}

func (m *MainWorkspace) FindWorkspace() (string, error) {
	path, err := os.Getwd()
	if err != nil {
		return "", err
	}

	volume := filepath.VolumeName(path)
	sentinel_filenames := []string{"WORKSPACE.bzlmod", "WORKSPACE.bazel", "MODULE.bazel", "WORKSPACE"} // search order

	for {
		// filepath.Dir() includes a trailing separator if we're at the root
		if path == volume+string(filepath.Separator) {
			path = volume
		}

		for _, sentinel := range sentinel_filenames {
			// Check if we're at the workspace path
			if s, err := os.Stat(filepath.Join(path, sentinel)); err == nil {
				// In macOS directories called "workspace" will match "WORKSPACE"
				// because the file system isn't case sensitive

				if !s.IsDir() && s.Name() == sentinel {
					return path, nil
				}
			}
		}

		// If we've reached the root, then we know the cwd isn't within a workspace
		if path == volume {
			return "", errors.New("ibazel was not invoked from within a workspace\n")
		}

		// Move one level up the path
		path = filepath.Dir(path)
	}
}

func (m *MainWorkspace) ExecuteCommand(command string, args []string) {
	for i, arg := range args {
		args[i] = strings.TrimSpace(arg)
	}
	log.Logf("Executing command: %s", command)
	workspacePath, err := m.FindWorkspace()
	if err != nil {
		log.Fatalf("Error finding workspace: %v", err)
		os.Exit(5)
	}
	log.Logf("Workspace path: %s", workspacePath)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, command, args...)
	log.Logf("Executing command: `%s`", strings.Join(cmd.Args, " "))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = workspacePath

	err = cmd.Run()
	if err != nil {
		log.Errorf("Command failed: %s %s. Error: %s", command, args, err)
	}
}

type FakeWorkspace struct{}

func (f *FakeWorkspace) FindWorkspace() (string, error) {
	return "", nil
}

func (f *FakeWorkspace) ExecuteCommand(command string, args []string) {}
