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
	"reflect"
	"testing"
)

func TestParsingArgs(t *testing.T) {
	for _, c := range []struct {
		in        []string
		targets   []string
		bazelArgs []string
		args      []string
		debugArgs [][]string
	}{
		// Empty case.
		{[]string{}, nil, nil, nil, nil},
		// Only one target.
		{[]string{"//my/target"}, []string{"//my/target"}, nil, nil, [][]string{{}}},
		// Only targets.
		{[]string{"//my/target1", "//my/target2"}, []string{"//my/target1", "//my/target2"}, nil, nil, [][]string{{},{}}},
		// arguments after a --.
		{[]string{"--", "--my_program_flag"}, nil, nil, []string{"--my_program_flag"}, nil},
		// Whitelisted bazel flag.
		{[]string{"--test_output=streaming"}, nil, []string{"--test_output=streaming"}, nil, nil},
		// Whitelisted bazel flag, arg, and target.
		{[]string{"--test_output=streaming", "--", "--my_program_flag"}, nil, []string{"--test_output=streaming"}, []string{"--my_program_flag"}, nil},
		// Multiple targets with multiple arguments.
		{[]string{"//my/target1", "--arg=t1_arg1", "--arg=t1_arg2", "//my/target2"}, []string{"//my/target1", "//my/target2"}, nil, nil, [][]string{{"t1_arg1", "t1_arg2"},{}}},
		// Multiple targets with single argument.
		{[]string{"//my/target1", "--arg=t1_arg1", "//my/target2", "--arg=t2_arg1", "//my/target3", "--arg=t3_arg1"}, []string{"//my/target1", "//my/target2", "//my/target3"}, nil, nil, [][]string{{"t1_arg1"}, {"t2_arg1"}, {"t3_arg1"}}},
		// Multiple targets with argument with whitespace.
		{[]string{"//my/target1", "--arg=t1 arg1", "//my/target2", "--arg=t2 arg1", "//my/target3"}, []string{"//my/target1", "//my/target2", "//my/target3"}, nil, nil, [][]string{{"\"t1 arg1\""}, {"\"t2 arg1\""}, {}}},
	} {
		targets, bazelArgs, args, debugArgs := parseArgs(c.in)
		if !reflect.DeepEqual(c.targets, targets) {
			t.Errorf("Targets not equal for args: %v\nGot:  %v\nWant: %v",
				c.in, targets, c.targets)
		}
		if !reflect.DeepEqual(c.bazelArgs, bazelArgs) {
			t.Errorf("Bazel args not equal for args: %v\nGot:  %v\nWant: %v",
				c.in, bazelArgs, c.bazelArgs)
		}
		if !reflect.DeepEqual(c.args, args) {
			t.Errorf("Additional args not equal for args: %v\nGot:  %v\nWant: %v",
				c.in, args, c.args)
		}
		if !reflect.DeepEqual(c.debugArgs, debugArgs) {
			t.Errorf("Debug args not equal for debugArgs: %v\nGot:  %v\nWant: %v",
				c.in, debugArgs, c.debugArgs)
		}
	}
}

func TestIsOverrideableBazelFlag(t *testing.T) {
	// Set some extra flags for testing
	overrideableBazelFlags = []string{
		"--test_output=",
	}

	for _, c := range []struct {
		arg          string
		overrideable bool
	}{
		{"--test", false},
		{"--i_love_ponies", false},
		{"--test_output", false},
		{"--test_output=streamed", true},
		{"--test_output_with_mooses=false", false},
	} {
		if isOverrideableBazelFlag(c.arg) != c.overrideable {
			t.Errorf("isOverrideableBazelFlag(%v) == %v (Got), Wanted %v", c.arg, !c.overrideable, c.overrideable)
		}
	}
}
