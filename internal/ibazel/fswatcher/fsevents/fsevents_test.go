// Copyright 2022 The Bazel Authors. All rights reserved.
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

//go:build darwin
// +build darwin

package fsevents

import (
	"github.com/google/go-cmp/cmp"
	"testing"
)

func TestFindCommonRoot(t *testing.T) {
	tests := []struct {
		in   []string
		want []string
	}{
		// Finds common sub-root of two directories
		{
			[]string{
				"/a/b/c/",
				"/a/d/",
			},
			[]string{"/a/"},
		},
		// Finds common sub-root of multiple recursive directories
		{
			[]string{
				"/a/b/c",
				"/a/b/c/e",
				"/a/b/d/e",
				"/a/b/d/f",
			},
			[]string{"/a/b/"},
		},
		// Returns the single input.
		{
			[]string{
				"/a/b/c/",
			},
			[]string{"/a/b/c/"},
		},
		// Returns an empty slice if there are no inputs.
		{
			[]string{},
			[]string{},
		},
	}
	for _, test := range tests {
		got, err := findCommonRoot(test.in)
		if err != nil {
			t.Errorf("unexpected error %v", err.Error())
		}
		if diff := cmp.Diff(got, test.want); diff != "" {
			t.Errorf("findCommonRoot diff (-got,+want):\n%s", diff)
		}
	}
}

func TestNoCommonRootError(t *testing.T) {
	_, err := findCommonRoot([]string{"/a/", "/b/"})
	if err == nil {
		t.Error("expected error when there is no common root")
	}
}
