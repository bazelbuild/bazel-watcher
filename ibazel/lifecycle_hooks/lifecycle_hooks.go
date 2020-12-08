// Copyright 2018 The Bazel Authors. All rights reserved.
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

// Implements lifecycle hooks for each command execution iteration.

package lifecycle_hooks

import (
	"bytes"
	"flag"

	"github.com/bazelbuild/bazel-watcher/ibazel/log"
	"github.com/bazelbuild/bazel-watcher/ibazel/workspace"
	"github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf/blaze_query"
	"github.com/mattn/go-shellwords"
)

var (
	runCommandBefore       = flag.String("run_command_before", "", "A command to run before each execution")
	runCommandAfter        = flag.String("run_command_after", "", "A command to run after each execution")
	runCommandAfterSuccess = flag.String("run_command_after_success", "", "A command to run after each successful execution")
)

type LifecycleHooks struct {
	w workspace.Workspace
}

func New() *LifecycleHooks {
	return &LifecycleHooks{
		w: &workspace.MainWorkspace{},
	}
}

func (l *LifecycleHooks) Initialize(info *map[string]string) {}

func (l *LifecycleHooks) TargetDecider(rule *blaze_query.Rule) {}

func (l *LifecycleHooks) ChangeDetected(targets []string, changeType string, change string) {}

func (l *LifecycleHooks) Cleanup() {}

func (l *LifecycleHooks) BeforeCommand(targets []string, command string) {
	l.parseAndExecuteCommand(*runCommandBefore)
}

func (l *LifecycleHooks) AfterCommand(targets []string, command string, success bool, output *bytes.Buffer) {
	l.parseAndExecuteCommand(*runCommandAfter)
	if success {
		l.parseAndExecuteCommand(*runCommandAfterSuccess)
	}
}

func (l *LifecycleHooks) parseAndExecuteCommand(commandToRun string) {
	if commandToRun != "" {
		commandAndArgs, err := shellwords.Parse(commandToRun)

		if err != nil {
			log.Fatalf("Fail to run command `%s`: %v", commandToRun, err)
		}

		l.w.ExecuteCommand(commandAndArgs[0], commandAndArgs[1:])
	}
}
