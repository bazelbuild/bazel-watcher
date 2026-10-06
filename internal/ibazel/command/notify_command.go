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

package command

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/bazelbuild/bazel-watcher/internal/ibazel/bep"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/log"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/process_group"
)

type notifyCommand struct {
	target       string
	startupArgs  []string
	bazelArgs    []string
	args         []string
	pg           process_group.ProcessGroup
	stdin        io.WriteCloser
	structured   bool
	outputGroups []string
	termSync     sync.Once
}

type buildEvent struct {
	Version                 int                     `json:"version"`
	Type                    string                  `json:"type"`
	Success                 bool                    `json:"success"`
	Changes                 []Change                `json:"changes"`
	AffectedTargets         []string                `json:"affected_targets,omitempty"`
	AffectedTargetsComplete bool                    `json:"affected_targets_complete,omitempty"`
	OutputGroups            map[string][]bep.Output `json:"output_groups,omitempty"`
	OutputGroupsComplete    bool                    `json:"output_groups_complete,omitempty"`
}

// NotifyCommand is an alternate mode for starting a command. In this mode the
// command will be notified on stdin that the source files have changed.
func NotifyCommand(startupArgs []string, bazelArgs []string, target string, args []string, structured bool, outputGroups []string) Command {
	return &notifyCommand{
		startupArgs:  startupArgs,
		target:       target,
		bazelArgs:    bazelArgs,
		args:         args,
		structured:   structured,
		outputGroups: outputGroups,
	}
}

func (c *notifyCommand) Terminate() {
	if !c.IsSubprocessRunning() {
		c.pg = nil
		return
	}
	c.termSync.Do(func() {
		terminate(c.pg)
	})
	c.pg = nil
}

func (c *notifyCommand) Kill() {
	if c.pg != nil {
		kill(c.pg)
	}
}

func (c *notifyCommand) Start() (*bytes.Buffer, error) {
	b := bazelNew()
	b.SetStartupArgs(c.startupArgs)
	b.SetArguments(c.argumentsWithOutputGroups())

	b.WriteToStderr(true)
	b.WriteToStdout(true)

	var outputBuffer *bytes.Buffer
	outputBuffer, c.pg = start(b, c.target, c.args)
	// Keep the writer around.
	var err error
	c.stdin, err = c.pg.RootProcess().StdinPipe()
	if err != nil {
		log.Errorf("Error getting stdin pipe: %v", err)
		return outputBuffer, err
	}

	c.pg.RootProcess().Env = append(os.Environ(), "IBAZEL_NOTIFY_CHANGES=y")

	if err = c.pg.Start(); err != nil {
		log.Errorf("Error starting process: %v", err)
		return outputBuffer, err
	}
	log.Log("Starting...")
	c.termSync = sync.Once{}
	return outputBuffer, nil
}

func (c *notifyCommand) NotifyOfChanges(changes []Change, impact ChangeImpact) *bytes.Buffer {
	b := bazelNew()
	b.SetStartupArgs(c.startupArgs)
	bepPath, cleanup := c.buildEventFile()
	defer cleanup()
	bazelArgs := c.argumentsWithOutputGroups()
	if bepPath != "" {
		bazelArgs = append(bazelArgs, "--build_event_json_file="+bepPath)
	}
	b.SetArguments(bazelArgs)

	b.WriteToStderr(true)
	b.WriteToStdout(true)

	_, err := c.stdin.Write([]byte("IBAZEL_BUILD_STARTED\n"))
	if err != nil {
		log.Errorf("Error writing build to stdin: %s", err)
	}

	outputBuffer, res := b.Norun(c.target)
	if res != nil {
		log.Errorf("IBAZEL BUILD FAILURE: %v", res)
		_, err := c.stdin.Write([]byte("IBAZEL_BUILD_COMPLETED FAILURE\n"))
		if err != nil {
			log.Errorf("Error writing failure to stdin: %s", err)
		}
		c.writeBuildEvent(false, changes, impact, nil, false)
	} else {
		log.Log("IBAZEL BUILD SUCCESS")
		_, err := c.stdin.Write([]byte("IBAZEL_BUILD_COMPLETED SUCCESS\n"))
		if err != nil {
			log.Errorf("Error writing success to stdin: %v", err)
		}
		outputGroups, complete := c.readOutputGroups(bepPath)
		c.writeBuildEvent(true, changes, impact, outputGroups, complete)
		if !c.IsSubprocessRunning() {
			log.Log("Restarting process...")
			c.Terminate()
			c.Start()
		}
	}
	return outputBuffer
}

func (c *notifyCommand) writeBuildEvent(success bool, changes []Change, impact ChangeImpact, outputGroups map[string][]bep.Output, outputGroupsComplete bool) {
	if !c.structured {
		return
	}
	event, err := json.Marshal(buildEvent{
		Version:                 1,
		Type:                    "build_completed",
		Success:                 success,
		Changes:                 changes,
		AffectedTargets:         impact.Targets,
		AffectedTargetsComplete: impact.Complete,
		OutputGroups:            outputGroups,
		OutputGroupsComplete:    outputGroupsComplete,
	})
	if err != nil {
		log.Errorf("Error encoding build event: %v", err)
		return
	}
	if _, err := c.stdin.Write(append(append([]byte("IBAZEL_EVENT "), event...), '\n')); err != nil {
		log.Errorf("Error writing build event to stdin: %v", err)
	}
}

func (c *notifyCommand) argumentsWithOutputGroups() []string {
	args := append([]string(nil), c.bazelArgs...)
	if len(c.outputGroups) == 0 {
		return args
	}
	groups := make([]string, 0, len(c.outputGroups))
	for _, group := range c.outputGroups {
		groups = append(groups, "+"+group)
	}
	return append(args, "--output_groups="+strings.Join(groups, ","))
}

func (c *notifyCommand) buildEventFile() (string, func()) {
	if !c.structured || len(c.outputGroups) == 0 {
		return "", func() {}
	}
	file, err := os.CreateTemp("", "ibazel-bep-*.json")
	if err != nil {
		log.Errorf("Error creating build event file: %v", err)
		return "", func() {}
	}
	if err := file.Close(); err != nil {
		log.Errorf("Error closing build event file: %v", err)
		_ = os.Remove(file.Name())
		return "", func() {}
	}
	return file.Name(), func() { _ = os.Remove(file.Name()) }
}

func (c *notifyCommand) readOutputGroups(path string) (map[string][]bep.Output, bool) {
	if path == "" {
		return nil, false
	}
	file, err := os.Open(path)
	if err != nil {
		log.Errorf("Error opening build events: %v", err)
		return nil, false
	}
	defer file.Close()
	groups, err := bep.ReadOutputGroups(file, c.outputGroups)
	if err != nil {
		log.Errorf("Error reading build event output groups: %v", err)
		return nil, false
	}
	return groups, true
}

func (c *notifyCommand) IsSubprocessRunning() bool {
	return c.pg != nil && subprocessRunning(c.pg.RootProcess())
}
