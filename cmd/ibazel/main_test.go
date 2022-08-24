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
		in          []string
		targets     []string
		startupArgs []string
		bazelArgs   []string
		args        []string
	}{
		// Empty case.
		{[]string{}, nil, nil, nil, nil},
		// Only targets.
		{[]string{"//my/target"}, []string{"//my/target"}, nil, nil, nil},
		// arguments after a --.
		{[]string{"--", "--my_program_flag"}, nil, nil, nil, []string{"--my_program_flag"}},
		// Whitelisted startup argument.
		{[]string{"--bazelrc=/home/libsamek/bazelrc", "--nohome_rc"}, nil, []string{"--bazelrc=/home/libsamek/bazelrc", "--nohome_rc"}, nil, nil},
		// Whitelisted bazel flag.
		{[]string{"--test_output=streaming"}, nil, nil, []string{"--test_output=streaming"}, nil},
		// Whitelisted bazel flag, arg, and target.
		{[]string{"--test_output=streaming", "--", "--my_program_flag"}, nil, nil,[]string{"--test_output=streaming"}, []string{"--my_program_flag"}},
	} {
		targets, startupArgs, bazelArgs, args := parseArgs(c.in)
		if !reflect.DeepEqual(c.targets, targets) {
			t.Errorf("Targets not equal for args: %v\nGot:  %v\nWant: %v",
				c.in, targets, c.targets)
		}
		if !reflect.DeepEqual(c.startupArgs, startupArgs) {
			t.Errorf("Startup arguments not equal for args: %v\nGot:  %v\nWant: %v",
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
