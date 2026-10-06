// Copyright 2026 The Bazel Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bep

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadOutputGroups(t *testing.T) {
	stream := strings.NewReader(`
{"id":{"namedSet":{"id":"child"}},"namedSetOfFiles":{"files":[{"name":"schema.ts","uri":"file:///execroot/bazel-out/bin/schema.ts","pathPrefix":["bazel-out","bin"],"digest":"schema-digest","length":"42"}]}}
{"id":{"namedSet":{"id":"root"}},"namedSetOfFiles":{"files":[{"name":"routes.ts","uri":"file:///execroot/bazel-out/bin/routes.ts","pathPrefix":["bazel-out","bin"],"digest":"routes-digest"}],"fileSets":[{"id":"child"}]}}
{"id":{"targetCompleted":{"label":"//app:dev"}},"completed":{"success":true,"outputGroup":[{"name":"generated","fileSets":[{"id":"root"}],"inlineFiles":[{"name":"metadata.json","pathPrefix":["bazel-out","bin"],"contents":"e30=","length":"2"},{"name":"link","pathPrefix":["bazel-out","bin"],"symlinkTargetPath":"schema.ts"}]},{"name":"manifest","fileSets":[]},{"name":"default","fileSets":[]}]}}
`)

	got, err := ReadOutputGroups(stream, []string{"generated", "manifest"})
	if err != nil {
		t.Fatal(err)
	}
	want := map[string][]Output{
		"generated": {
			{Path: "bazel-out/bin/link", SymlinkTargetPath: "schema.ts"},
			{Path: "bazel-out/bin/metadata.json", Contents: "e30=", Length: 2},
			{Path: "bazel-out/bin/routes.ts", URI: "file:///execroot/bazel-out/bin/routes.ts", Digest: "routes-digest"},
			{Path: "bazel-out/bin/schema.ts", URI: "file:///execroot/bazel-out/bin/schema.ts", Digest: "schema-digest", Length: 42},
		},
		"manifest": {},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("output groups diff (-want +got):\n%s", diff)
	}
}

func TestReadOutputGroupsRejectsMissingGroup(t *testing.T) {
	stream := strings.NewReader(`{"id":{"targetCompleted":{"label":"//app:dev"}},"completed":{"success":true,"outputGroup":[{"name":"default","fileSets":[]}]}}`)

	_, err := ReadOutputGroups(stream, []string{"generated"})
	if err == nil || !strings.Contains(err.Error(), `output group "generated" was not reported`) {
		t.Fatalf("ReadOutputGroups() error = %v, want missing group error", err)
	}
}

func TestReadOutputGroupsRejectsIncompleteGroup(t *testing.T) {
	stream := strings.NewReader(`{"id":{"targetCompleted":{"label":"//app:dev"}},"completed":{"success":true,"outputGroup":[{"name":"generated","incomplete":true}]}}`)

	_, err := ReadOutputGroups(stream, []string{"generated"})
	if err == nil || !strings.Contains(err.Error(), `output group "generated" is incomplete`) {
		t.Fatalf("ReadOutputGroups() error = %v, want incomplete group error", err)
	}
}

func TestReadOutputGroupsRejectsMissingNamedSet(t *testing.T) {
	stream := strings.NewReader(`{"id":{"targetCompleted":{"label":"//app:dev"}},"completed":{"success":true,"outputGroup":[{"name":"generated","fileSets":[{"id":"missing"}]}]}}`)

	_, err := ReadOutputGroups(stream, []string{"generated"})
	if err == nil || !strings.Contains(err.Error(), `missing named set "missing"`) {
		t.Fatalf("ReadOutputGroups() error = %v, want missing named set error", err)
	}
}

func TestReadOutputGroupsRejectsConflictingMetadata(t *testing.T) {
	stream := strings.NewReader(`
{"id":{"namedSet":{"id":"outputs"}},"namedSetOfFiles":{"files":[{"name":"result.txt","pathPrefix":["bazel-out","bin"],"digest":"one"},{"name":"result.txt","pathPrefix":["bazel-out","bin"],"digest":"two"}]}}
{"id":{"targetCompleted":{"label":"//app:dev"}},"completed":{"success":true,"outputGroup":[{"name":"generated","fileSets":[{"id":"outputs"}]}]}}
`)

	_, err := ReadOutputGroups(stream, []string{"generated"})
	if err == nil || !strings.Contains(err.Error(), `artifact "bazel-out/bin/result.txt" has conflicting metadata`) {
		t.Fatalf("ReadOutputGroups() error = %v, want conflicting metadata error", err)
	}
}

func TestReadOutputGroupsDeduplicatesNamedAndInlineFiles(t *testing.T) {
	stream := strings.NewReader(`
{"id":{"namedSet":{"id":"child"}},"namedSetOfFiles":{"files":[{"name":"result.txt","digest":"one"}]}}
{"id":{"namedSet":{"id":"root"}},"namedSetOfFiles":{"files":[{"name":"result.txt","digest":"one"}],"fileSets":[{"id":"child"},{"id":"root"}]}}
{"completed":{"outputGroup":[{"name":"generated","fileSets":[{"id":"root"},{"id":"child"}],"inlineFiles":[{"name":"result.txt","digest":"one"}]}]}}
`)
	got, err := ReadOutputGroups(stream, []string{"generated"})
	if err != nil {
		t.Fatal(err)
	}
	want := map[string][]Output{"generated": {{Path: "result.txt", Digest: "one"}}}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("output groups diff (-want +got):\n%s", diff)
	}
}

func TestReadOutputGroupsRejectsConflictingInlineMetadata(t *testing.T) {
	stream := strings.NewReader(`
{"id":{"namedSet":{"id":"outputs"}},"namedSetOfFiles":{"files":[{"name":"result.txt","digest":"one"}]}}
{"completed":{"outputGroup":[{"name":"generated","fileSets":[{"id":"outputs"}],"inlineFiles":[{"name":"result.txt","digest":"two"}]}]}}
`)
	_, err := ReadOutputGroups(stream, []string{"generated"})
	if err == nil || !strings.Contains(err.Error(), `artifact "result.txt" has conflicting metadata`) {
		t.Fatalf("ReadOutputGroups() error = %v, want conflicting metadata error", err)
	}
}
