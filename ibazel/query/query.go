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

package query

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bazelbuild/bazel-watcher/bazel"
	blaze_query "github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf"
)

type BazelQuerier struct {
	bazelNew func() bazel.Bazel

	workspacePath string
}

func New(bazelNew func() bazel.Bazel, workspacePath string) *BazelQuerier {
	return &BazelQuerier{
		bazelNew:      bazelNew,
		workspacePath: workspacePath,
	}
}

func (bq *BazelQuerier) QueryForSourceFiles(query string) ([]string, error) {
	b := bq.bazelNew()

	res, err := b.Query(query)
	if err != nil {
		return nil, err
	}

	toWatch := make([]string, 0, 10000)
	for _, target := range res.Target {
		switch *target.Type {
		case blaze_query.Target_SOURCE_FILE:
			label := *target.SourceFile.Name
			if strings.HasPrefix(label, "@") {
				continue
			}
			if strings.HasPrefix(label, "//external") {
				continue
			}

			// For files that are served from the root they will being with "//:". This
			// is a problematic string because, for example, "//:demo.sh" will become
			// "/demo.sh" which is in the root of the filesystem and is unlikely to exist.
			if strings.HasPrefix(label, "//:") {
				label = label[3:]
			}

			label = strings.Replace(strings.TrimPrefix(label, "//"), ":", string(filepath.Separator), 1)
			toWatch = append(toWatch, filepath.Join(bq.workspacePath, label))
			break
		default:
			fmt.Fprintf(os.Stderr, "%v\n\n", target)
		}
	}

	return toWatch, nil
}

func (bq *BazelQuerier) QueryRule(rule string) (*blaze_query.Rule, error) {
	b := bq.bazelNew()

	res, err := b.Query(rule)
	if err != nil {
		return nil, err
	}

	for _, target := range res.Target {
		switch *target.Type {
		case blaze_query.Target_RULE:
			return target.Rule, nil
		}
	}

	return nil, errors.New("No information available")
}
