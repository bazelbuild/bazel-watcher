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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bazelbuild/bazel-watcher/internal/ibazel/log"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/workspace"
	"github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf/blaze_query"
)

var (
	runOutput = flag.Bool(
		"run_output",
		true,
		"Search for commands in Bazel output that match a regex and execute them, the default path of file should be in the workspace root .bazel_fix_commands.json")
	runOutputInteractive = flag.Bool(
		"run_output_interactive",
		true,
		"Use an interactive prompt when executing commands in Bazel output")
	notifiedUser = false
)

// This RegExp will match ANSI escape codes.
var escapeCodeCleanerRegex = regexp.MustCompile("\\x1B\\[[\\x30-\\x3F]*[\\x20-\\x2F]*[\\x40-\\x7E]")

type OutputRunner struct {
	w workspace.Workspace
}

type Optcmd struct {
	Regex   string   `json:"regex"`
	Command string   `json:"command"`
	Args    []string `json:"args"`
}

func New() *OutputRunner {
	i := &OutputRunner{
		w: &workspace.MainWorkspace{},
	}
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

	optcmd := i.readConfigs(jsonCommandPath)
	if optcmd == nil {
		optcmd = []Optcmd{defaultRegex}
		log.Log("using default bazel_fix_commands")
	}
	log.Logf("sifting output: %s", output.String())
	commandLines, commands, args := matchRegex(optcmd, output)
	log.Logf("using %d matching commandlines", len(commandLines))
	for idx, _ := range commandLines {
		log.Logf("matched command number %d: %v", idx, commands[idx])
		if *runOutputInteractive {
			if i.promptCommand(commands[idx], args[idx]) {
				i.executeCommand(commands[idx], args[idx])
			}
		} else {
			i.executeCommand(commands[idx], args[idx])
		}
	}
}

func (o *OutputRunner) readConfigs(configPath string) []Optcmd {
	workspacePath, err := o.w.FindWorkspace()
	if err != nil {
		log.Fatalf("Error finding workspace: %v", err)
		os.Exit(5)
	}

	jsonFile, err := os.Open(filepath.Join(workspacePath, configPath))
	if os.IsNotExist(err) {
		// Note this is not attached to the os.IsNotExist because we don't want the
		// other error handler to catch if we hav already notified.
		if !notifiedUser {
			log.Banner(
				"Did you know iBazel can invoke programs like Gazelle, buildozer, and",
				"other BUILD file generators for you automatically based on bazel output?",
				"Documentation at: https://github.com/bazelbuild/bazel-watcher#output-runner")
		}
		notifiedUser = true
		return nil
	} else if err != nil {
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
	distinctCommands := map[string]bool{}
	scanner := bufio.NewScanner(output)
	for scanner.Scan() {
		line := escapeCodeCleanerRegex.ReplaceAllLiteralString(scanner.Text(), "")
		for _, oc := range optcmd {
			re := regexp.MustCompile(oc.Regex)
			matches := re.FindStringSubmatch(line)
			if matches != nil && len(matches) >= 0 {
				command := convertArg(matches, oc.Command)
				cmdArgs := convertArgs(matches, oc.Args)
				fullCmd := strings.Join(append([]string{command}, cmdArgs...), " ")
				if _, found := distinctCommands[fullCmd]; !found {
					commandLines = append(commandLines, matches[0])
					commands = append(commands, command)
					args = append(args, cmdArgs)
					distinctCommands[fullCmd] = true
				}
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
		var converted strings.Builder
		converted.Grow(len(args[i]))

		matchIndex := 0
		matching := false

		writeMatch := func() {
			if matching {
				if matchIndex < len(matches) {
					converted.WriteString(matches[matchIndex])
				} else {
					converted.WriteRune('$')
					converted.WriteString(strconv.Itoa(matchIndex))
				}
				matchIndex = 0
				matching = false
			}
		}

		for _, c := range args[i] {
			if c == '$' {
				if matching {
					converted.WriteRune(c)
				}
				matching = !matching
			} else if matching && c >= '0' && c <= '9' {
				matchIndex = matchIndex*10 + int(c-'0')
			} else {
				writeMatch()
				converted.WriteRune(c)
			}
		}
		writeMatch()

		rst = append(rst, converted.String())
	}
	return rst
}

func (_ *OutputRunner) promptCommand(command string, args []string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fprintf(os.Stderr, "Do you want to execute this command?\n%s %s\n[y/N]", command, strings.Join(args, " "))
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

func (o *OutputRunner) executeCommand(command string, args []string) {
	o.w.ExecuteCommand(command, args)
}

func (i *OutputRunner) Cleanup() {}
