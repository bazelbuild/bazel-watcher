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
	"fmt"
	"os/exec"
	"regexp"
	"testing"

	"github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf/analysis"
	"github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf/blaze_query"
	"github.com/google/go-cmp/cmp"
)

type MockBazel struct {
	actions        [][]string
	queryResponse  map[string]*blaze_query.QueryResult
	cqueryResponse map[string]*analysis.CqueryResult
	args           []string
	startupArgs    []string
	info           map[string]string

	buildError error
	waitError  error
}

func (b *MockBazel) Args() []string {
	return b.args
}

func (b *MockBazel) SetArguments(args []string) {
	b.actions = append(b.actions, append([]string{"SetArguments"}, args...))
	b.args = args
}

func (b *MockBazel) SetStartupArgs(args []string) {
	b.actions = append(b.actions, append([]string{"SetStartupArgs"}, args...))
	b.startupArgs = args
}

func (b *MockBazel) SetInfo(info map[string]string) {
	b.info = info
}
func (b *MockBazel) WriteToStderr(v bool) {
	b.actions = append(b.actions, []string{"WriteToStderr", fmt.Sprint(v)})
}
func (b *MockBazel) WriteToStdout(v bool) {
	b.actions = append(b.actions, []string{"WriteToStdout", fmt.Sprint(v)})
}
func (b *MockBazel) Info() (map[string]string, error) {
	b.actions = append(b.actions, []string{"Info"})
	return b.info, nil
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

	if !ok {
		var candidates []string
		for candidate := range b.queryResponse {
			candidates = append(candidates, candidate)
		}
		panic(fmt.Sprintf("Unable to find query result for %q. Only have %v.", query, candidates))
	}

	return res, nil
}
func (b *MockBazel) AddCQueryResponse(query string, res *analysis.CqueryResult) {
	if b.cqueryResponse == nil {
		b.cqueryResponse = map[string]*analysis.CqueryResult{}
	}
	b.cqueryResponse[query] = res
}
func (b *MockBazel) CQuery(args ...string) (*analysis.CqueryResult, error) {
	b.actions = append(b.actions, append([]string{"CQuery"}, args...))
	query := args[0]
	res, ok := b.cqueryResponse[query]

	if !ok {
		var candidates []string
		for candidate := range b.cqueryResponse {
			candidates = append(candidates, candidate)
		}
		panic(fmt.Sprintf("Unable to find cquery result for %q. Only have %v.", query, candidates))
	}

	return res, nil
}
func (b *MockBazel) Build(args ...string) (*bytes.Buffer, error) {
	b.actions = append(b.actions, append([]string{"Build"}, args...))
	return nil, b.buildError
}
func (b *MockBazel) BuildError(e error) {
	b.buildError = e
}
func (b *MockBazel) Test(args ...string) (*bytes.Buffer, error) {
	b.actions = append(b.actions, append([]string{"Test"}, args...))
	return nil, nil
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

	if diff := cmp.Diff(b.actions, expected, cmp.FilterValues(func(a, b string) bool {
		return true
	}, cmp.Comparer(func(a, b string) bool {
		{
			match, _ := regexp.MatchString(a, b)
			if match {
				return true
			}
		}
		{
			match, _ := regexp.MatchString(b, a)
			if match {
				return true
			}
		}
		return a == b
	}))); diff != "" {
		t.Errorf("Action diff (-got (%d),+want (%d)):\n%s", len(b.actions), len(expected), diff)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
