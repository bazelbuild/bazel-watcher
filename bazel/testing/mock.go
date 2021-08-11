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

package testing

import (
	"bytes"
	"os/exec"
	"regexp"
	"testing"

	"github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf/analysis"
	"github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf/blaze_query"
)

type MockBazel struct {
	actions        [][]string
	queryResponse  map[string]*blaze_query.QueryResult
	cqueryResponse map[string]*analysis.CqueryResult
	args           []string
	startupArgs    []string

	buildError error
	waitError  error
}

func (b *MockBazel) SetArguments(args []string) {
	b.args = args
}

func (b *MockBazel) SetStartupArgs(args []string) {
	b.startupArgs = args
}

func (b *MockBazel) WriteToStderr(v bool) {
	b.actions = append(b.actions, []string{"WriteToStderr"})
}
func (b *MockBazel) WriteToStdout(v bool) {
	b.actions = append(b.actions, []string{"WriteToStdout"})
}
func (b *MockBazel) Info() (map[string]string, error) {
	b.actions = append(b.actions, []string{"Info"})
	return map[string]string{}, nil
}
func (b *MockBazel) AddQueryResponse(query string, res *blaze_query.QueryResult) {
	if b.queryResponse == nil {
		b.queryResponse = map[string]*blaze_query.QueryResult{}
	}
	b.queryResponse[query] = res
}
func (b *MockBazel) Query(args ...string) (*blaze_query.QueryResult, error) {
	b.actions = append(b.actions, append([]string{"Query"}, args...))
	query := args[0]
	res, ok := b.queryResponse[query]

	if !ok || res == nil {
		res = &blaze_query.QueryResult{}
	}

	return res, nil
}
func (b *MockBazel) AddCQueryResponse(query string, res *analysis.CqueryResult) {
	if b.queryResponse == nil {
		b.cqueryResponse = map[string]*analysis.CqueryResult{}
	}
	b.cqueryResponse[query] = res
}
func (b *MockBazel) CQuery(args ...string) (*analysis.CqueryResult, error) {
	b.actions = append(b.actions, append([]string{"CQuery"}, args...))
	query := args[0]
	res, ok := b.cqueryResponse[query]

	if !ok || res == nil {
		res = &analysis.CqueryResult{}
	}

	return res, nil
}
func (b *MockBazel) Build(args ...string) (*bytes.Buffer, error) {
	b.actions = append(b.actions, append([]string{"Build"}, args...))
	return nil, b.buildError
}
func (b *MockBazel) BuildCancelable(cancelCh chan bool, args ...string) (*bytes.Buffer, error) {
	return b.Build(args...)
}
func (b *MockBazel) BuildError(e error) {
	b.buildError = e
}
func (b *MockBazel) Test(args ...string) (*bytes.Buffer, error) {
	b.actions = append(b.actions, append([]string{"Test"}, args...))
	return nil, nil
}
func (b *MockBazel) TestCancelable(cancelCh chan bool, args ...string) (*bytes.Buffer, error) {
	return b.Test(args...)
}
func (b *MockBazel) Run(args ...string) (*exec.Cmd, *bytes.Buffer, error) {
	b.actions = append(b.actions, append([]string{"Run"}, args...))
	return nil, nil, nil
}
func (b *MockBazel) WaitError(e error) {
	b.waitError = e
}
func (b *MockBazel) Wait() error {
	return b.waitError
}
func (b *MockBazel) Cancel() {
	b.actions = append(b.actions, []string{"Cancel"})
}
func (b *MockBazel) AssertActions(t *testing.T, expected [][]string) {
	t.Helper()

	if len(b.actions) != len(expected) {
		t.Errorf("Test didn't meet expectations len(b.actions) == %d != len(expected) == %d.\nWant: %#v\nGot:  %#v", len(b.actions), len(expected), expected, b.actions)
		return
	}

	for i := range b.actions {
		if len(b.actions[i]) != len(expected[i]) {
			t.Errorf("Test didn't meet expectations len(b.actions[%d]) == %d != len(expected[%d]) == %d.\nWant: %#v\nGot:  %#v", i, len(b.actions[i]), i, len(expected[i]), expected, b.actions)
			return
		}

		for j := range b.actions[i] {
			match, _ := regexp.MatchString(expected[i][j], b.actions[i][j])
			if !match {
				t.Errorf("Test didn't meet expectations.\nWant: %#v\nGot:  %#v", expected, b.actions)
				return
			}
		}
	}
}
