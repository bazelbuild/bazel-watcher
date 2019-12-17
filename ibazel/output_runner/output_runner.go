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

package output_runner

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/bazelbuild/bazel-watcher/ibazel/log"
	"github.com/bazelbuild/bazel-watcher/ibazel/workspace_finder"
	blaze_query "github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf"
)

var runOutput = flag.Bool(
	"run_output",
	false,
	"Search for commands in Bazel output that match a regex and execute them, the default path of file should be in the workspace root .bazel_fix_commands.json")
var runOutputInteractive = flag.Bool(
	"run_output_interactive",
	true,
	"Use an interactive prompt when executing commands in Bazel output")

// This RegExp will match ANSI escape codes.
var escapeCodeCleanerRegex = regexp.MustCompile("\\x1B\\[[\\x30-\\x3F]*[\\x20-\\x2F]*[\\x40-\\x7E]")

type OutputRunner struct{}

type Optcmd struct {
	Regex   string   `json:"regex"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func New() *OutputRunner {
	i := &OutputRunner{}
	return i
}

func (i *OutputRunner) Initialize(info *map[string]string) {}

func (i *OutputRunner) TargetDecider(rule *blaze_query.Rule) {}

func (i *OutputRunner) ChangeDetected(targets []string, changeType string, change string) {}

func (i *OutputRunner) BeforeCommand(targets []string, command string) {}

func (i *OutputRunner) AfterCommand(targets []string, command string, success bool, output *bytes.Buffer) {
	if *runOutput == false || output == nil {
		return
	}

	jsonCommandPath := ".bazel_fix_commands.json"
	defaultRegex := Optcmd{
		Regex:   "^buildozer '(.*)'\\s+(.*)$",
		Command: "buildozer",
		Args:    []string{"$1", "$2"},
	}

	optcmd := readConfigs(jsonCommandPath)
	if optcmd == nil {
		log.Log("Use default regex")
		optcmd = []Optcmd{defaultRegex}
	}
	commandLines, commands, args := matchRegex(optcmd, output)
	for idx, _ := range commandLines {
		if *runOutputInteractive {
			if promptCommand(commandLines[idx]) {
				executeCommand(commands[idx], args[idx])
			}
		} else {
			executeCommand(commands[idx], args[idx])
		}
	}
}

func readConfigs(configPath string) []Optcmd {
	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Errorf("Error reading config: %s", err)
		return nil
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var optcmd []Optcmd
	err = json.Unmarshal(byteValue, &optcmd)
	if err != nil {
		log.Errorf("Error in .bazel_fix_commands.json: %s", err)
	}

	return optcmd
}

func matchRegex(optcmd []Optcmd, output *bytes.Buffer) ([]string, []string, [][]string) {
	var commandLines, commands []string
	var args [][]string
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		line := escapeCodeCleanerRegex.ReplaceAllLiteralString(scanner.Text(), "")
		for _, oc := range optcmd {
			re := regexp.MustCompile(oc.Regex)
			matches := re.FindStringSubmatch(line)
			if matches != nil && len(matches) >= 0 {
				commandLines = append(commandLines, matches[0])
				commands = append(commands, convertArg(matches, oc.Command))
				args = append(args, convertArgs(matches, oc.Args))
			}
		}
	}
	return commandLines, commands, args
}

func convertArg(matches []string, arg string) string {
	if strings.HasPrefix(arg, "$") {
		val, _ := strconv.Atoi(arg[1:])
		return matches[val]
	}
	return arg
}

func convertArgs(matches []string, args []string) []string {
	var rst []string
	for i, _ := range args {
		if strings.HasPrefix(args[i], "$") {
			val, _ := strconv.Atoi(args[i][1:])
			rst = append(rst, matches[val])
		} else {
			rst = append(rst, args[i])
		}
	}
	return rst
}

func promptCommand(command string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stderr, "Do you want to execute this command?\n%s\n[y/N]", command)
	text, _ := reader.ReadString('\n')
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	text = strings.TrimRight(text, "\n")
	if text == "y" {
		return true
	} else {
		return false
	}
}

func executeCommand(command string, args []string) {
	for i, arg := range args {
		args[i] = strings.TrimSpace(arg)
	}
	log.Logf("Executing command: %s", command)
	workspaceFinder := &workspace_finder.MainWorkspaceFinder{}
	workspacePath, err := workspaceFinder.FindWorkspace()
	if err != nil {
		log.Fatalf("Error finding workspace: %v", err)
		os.Exit(5)
	}
	log.Logf("Workspace path: %s", workspacePath)

	ctx, _ := context.WithCancel(context.Background())
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

func (i *OutputRunner) Cleanup() {}
