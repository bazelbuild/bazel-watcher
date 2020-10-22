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
	"bytes"
	"reflect"
	"testing"

	"github.com/bazelbuild/bazel-watcher/ibazel/workspace"
)

func TestConvertArgs(t *testing.T) {
	matches := []string{"my_command my_arg1 my_arg2 my_arg3", "my_command", "my_arg1", "my_arg2", "my_arg3"}
	// Command parsing tests
	for _, c := range []struct {
		cmd   string
		truth string
	}{
		{"$1", "my_command"},
		{"warning", "warning"},
		{"keep_command", "keep_command"},
	} {
		new_cmd := convertArg(matches, c.cmd)
		if !reflect.DeepEqual(c.truth, new_cmd) {
			t.Errorf("Command not equal: %v\nGot:  %v\nWant: %v",
				c.cmd, new_cmd, c.truth)
		}
	}
	// Arguments parsing tests
	for _, c := range []struct {
		cmd   []string
		truth []string
	}{
		{[]string{"$2", "$3"}, []string{"my_arg1", "my_arg2"}},
		{[]string{"$2", "$3", "$4"}, []string{"my_arg1", "my_arg2", "my_arg3"}},
		{[]string{"$2", "dont_change_arg"}, []string{"my_arg1", "dont_change_arg"}},
		{[]string{"keep_arg", "$3"}, []string{"keep_arg", "my_arg2"}},
	} {
		new_cmd := convertArgs(matches, c.cmd)
		if !reflect.DeepEqual(c.truth, new_cmd) {
			t.Errorf("Command not equal: %v\nGot:  %v\nWant: %v",
				c.cmd, new_cmd, c.truth)
		}
	}
}

func TestReadConfigs(t *testing.T) {
	i := &OutputRunner{
		w: &workspace.FakeWorkspace{},
	}
	optcmd := i.readConfigs("output_runner_test.json")

	for idx, c := range []struct {
		regex   string
		command string
		args    []string
	}{
		{"^(buildozer) '(.*)'\\s+(.*)$", "$1", []string{"$2", "$3"}},
		{"WARNING", "warn", []string{"keep_calm", "dont_panic"}},
		{"DANGER", "danger", []string{"be_careful", "why_so_serious"}},
	} {
		if !reflect.DeepEqual(c.regex, optcmd[idx].Regex) {
			t.Errorf("Regex not equal: %v\nGot:  %v\nWant: %v",
				optcmd[idx], optcmd[idx].Regex, c.regex)
		}
		if !reflect.DeepEqual(c.command, optcmd[idx].Command) {
			t.Errorf("Command not equal: %v\nGot:  %v\nWant: %v",
				optcmd[idx], optcmd[idx].Command, c.command)
		}
		if !reflect.DeepEqual(c.args, optcmd[idx].Args) {
			t.Errorf("Args not equal: %v\nGot:  %v\nWant: %v",
				optcmd[idx], optcmd[idx].Args, c.args)
		}
	}
}

func TestMatchRegex(t *testing.T) {
	buf := bytes.Buffer{}
	buf.WriteString("buildozer 'add deps test_dep1' //target1:target1\n")
	buf.WriteString("buildozer 'add deps test_dep2' //target2:target2\n")
	buf.WriteString("Check that imports in Go sources match importpath attributes in deps.\n")
	buf.WriteString("Other build output which does not match things.\n")
	// Duplicate matches, to be deduplicated.
	buf.WriteString("Check that imports in Go sources match importpath attributes in deps.\n")
	buf.WriteString("buildozer 'add deps test_dep2' //target2:target2\n")

	optcmd := []Optcmd{
		{Regex: "^Check that imports in Go sources match importpath attributes in deps.$", Command: "bazel", Args: []string{"run", "//:gazelle"}},
		{Regex: "^(buildozer) '(.*)'\\s+(.*)$", Command: "$1", Args: []string{"$2", "$3"}},
	}

	_, commands, args := matchRegex(optcmd, &buf)

	expected := []struct {
		cls string
		cs  string
		as  []string
	}{
		{"buildozer 'add deps test_dep1' //target1:target1", "buildozer", []string{"add deps test_dep1", "//target1:target1"}},
		{"buildozer 'add deps test_dep2' //target2:target2", "buildozer", []string{"add deps test_dep2", "//target2:target2"}},
		{"Check that imports in Go sources match importpath attributes in deps.", "bazel", []string{"run", "//:gazelle"}},
	}
	if len(expected) != len(commands) {
		t.Errorf("Did not receive expected number of commands:\nGot: %d\nWant: %d", len(commands), len(expected))
	}
	for idx, c := range expected {
		if !reflect.DeepEqual(c.cs, commands[idx]) {
			t.Errorf("Commands not equal: %v\nGot:  %v\nWant: %v",
				c.cls, commands[idx], c.cs)
		}
		if !reflect.DeepEqual(c.as, args[idx]) {
			t.Errorf("Arguments not equal: %v\nGot:  %v\nWant: %v",
				c.cls, args[idx], c.as)
		}
	}
}

var cleanerTests = []struct {
	in  string
	out []string
}{
	{
		"buildozer 'add deps //wow' //fake:target0",
		[]string{"buildozer 'add deps //wow' //fake:target0"},
	},
	{
		"[96mbuildozer [0m[93m[0m[93m[0m[91m'add deps [0m[90m//wow'[0m //fake:target1",
		[]string{"buildozer 'add deps //wow' //fake:target1"},
	},
	{
		"build[96mozer 'add d[96meps //w[96mow' //fake:tar[96mget2",
		[]string{"buildozer 'add deps //wow' //fake:target2"},
	},
	{
		"[0m[90mbuildozer 'a[0m[93mdd deps //wow' [0m[93m//fake:target3[91m",
		[]string{"buildozer 'add deps //wow' //fake:target3"},
	},
	{
		"buildozer 'add deps //wow[0m[93m //fake:target4",
		[]string(nil),
	},
}

func TestMatchCleanRegex(t *testing.T) {
	optcmd := []Optcmd{
		{Regex: "^(buildozer) '(.*)'\\s+(.*)$", Command: "$1", Args: []string{"$2", "$3"}},
	}

	for _, tt := range cleanerTests {
		t.Run(tt.in, func(t *testing.T) {
			buf := bytes.Buffer{}
			buf.WriteString(tt.in)
			cmdLines, _, _ := matchRegex(optcmd, &buf)

			if !reflect.DeepEqual(cmdLines, tt.out) {
				t.Errorf("Commands not equal!\nGot:  %v\nWant: %v", cmdLines, tt.out)
			}
		})
	}

}
