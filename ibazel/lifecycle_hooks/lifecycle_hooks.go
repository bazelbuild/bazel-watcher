package lifecycle_hooks

import (
	"bytes"
	"flag"

	"github.com/bazelbuild/bazel-watcher/ibazel/workspace"
	"github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf/blaze_query"
	"github.com/mattn/go-shellwords"
)

var runCommandBefore = flag.String("run_command_before", "", "A command to run before each execution")

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
	if *runCommandBefore != "" {
		commandAndArgs, err := shellwords.Parse(*runCommandBefore)

		if err != nil {
			panic(err)
		}

		l.w.ExecuteCommand(commandAndArgs[0], commandAndArgs[1:])
	}
}

func (l *LifecycleHooks) AfterCommand(targets []string, command string, success bool, output *bytes.Buffer) {
}
