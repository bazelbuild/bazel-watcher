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

package bazel

import (
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	b := New()
	if b == nil {
		t.Fatalf("Created a nil object")
	}
}

func TestProcessInfo(t *testing.T) {
	b := &bazel{}
	got, err := b.processInfo(`KEY: VALUE
KEY2: VALUE2
KEY3: value`)
	if err != nil {
		t.Errorf("Error processing info", err)
	}

	expected := map[string]string{
		"KEY":  "VALUE",
		"KEY2": "VALUE2",
		"KEY3": "value",
	}

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Objects were unequal. Got:\n%s\nExpected:\n%s", got, expected)
	}
}

func TestWriteToStderrAndStdout(t *testing.T) {
	b := &bazel{}

	// By default it should write to its own pipe.
	b.newCommand("version")
	if b.cmd.Stdout == os.Stdout {
		t.Errorf("Set stdout to os.Stdout")
	}
	if b.cmd.Stderr == os.Stderr {
		t.Errorf("Set stderr to os.Stderr")
	}

	// If set to true it should write to the os version
	b.WriteToStderr(true)
	b.WriteToStdout(true)
	b.newCommand("version")
	if b.cmd.Stdout != os.Stdout {
		t.Errorf("Didn't set stdout to os.Stdout")
	}
	if b.cmd.Stderr != os.Stderr {
		t.Errorf("Didn't set stderr to os.Stderr")
	}

	// If set to false it should not write to the os version
	b.WriteToStderr(false)
	b.WriteToStdout(false)
	b.newCommand("version")
	if b.cmd.Stdout == os.Stdout {
		t.Errorf("Set stdout to os.Stdout")
	}
	if b.cmd.Stderr == os.Stderr {
		t.Errorf("Set stderr to os.Stderr")
	}
}

// Test that cancel doesn't NPE if there is no command running.
func TestCancel(t *testing.T) {
	b := New()
	b.Cancel()
}
