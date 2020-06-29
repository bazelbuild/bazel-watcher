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

package fswatcher

import (
	"github.com/bazelbuild/bazel-watcher/ibazel/fswatcher/common"
)

type Event = common.Event
type Op = common.Op

const Create = common.Create
const Write = common.Write
const Remove = common.Remove
const Rename = common.Rename
const Chmod = common.Chmod

type Watcher = common.Watcher
