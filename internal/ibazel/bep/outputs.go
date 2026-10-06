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

// Package bep reads the subset of Bazel's JSON Build Event Protocol needed by
// notification-mode commands.
package bep

import (
	"encoding/json"
	"fmt"
	"io"
	"path"
	"sort"
)

// Output identifies one artifact in a Bazel output group. Contents uses the
// base64 encoding from Bazel's JSON representation of the BEP.
type Output struct {
	Path              string `json:"path"`
	URI               string `json:"uri,omitempty"`
	Contents          string `json:"contents,omitempty"`
	SymlinkTargetPath string `json:"symlink_target_path,omitempty"`
	Digest            string `json:"digest,omitempty"`
	Length            int64  `json:"length,omitempty"`
}

type event struct {
	ID              eventID          `json:"id"`
	NamedSetOfFiles *namedSetOfFiles `json:"namedSetOfFiles"`
	Completed       *targetComplete  `json:"completed"`
}

type eventID struct {
	NamedSet *namedSetID `json:"namedSet"`
}

type namedSetID struct {
	ID string `json:"id"`
}

type namedSetOfFiles struct {
	Files    []file       `json:"files"`
	FileSets []namedSetID `json:"fileSets"`
}

type file struct {
	Name              string   `json:"name"`
	URI               string   `json:"uri"`
	Contents          string   `json:"contents"`
	SymlinkTargetPath string   `json:"symlinkTargetPath"`
	PathPrefix        []string `json:"pathPrefix"`
	Digest            string   `json:"digest"`
	Length            int64    `json:"length,string"`
}

type targetComplete struct {
	OutputGroups []outputGroup `json:"outputGroup"`
}

type outputGroup struct {
	Name        string       `json:"name"`
	FileSets    []namedSetID `json:"fileSets"`
	Incomplete  bool         `json:"incomplete"`
	InlineFiles []file       `json:"inlineFiles"`
}

// ReadOutputGroups extracts selected output groups from a JSON BEP stream.
func ReadOutputGroups(reader io.Reader, selected []string) (map[string][]Output, error) {
	wanted := make(map[string]struct{}, len(selected))
	type groupFiles struct {
		roots  []string
		inline []file
	}
	groups := make(map[string]*groupFiles, len(selected))
	seenGroups := make(map[string]struct{}, len(selected))
	for _, name := range selected {
		wanted[name] = struct{}{}
		groups[name] = &groupFiles{}
	}

	namedSets := make(map[string]namedSetOfFiles)
	decoder := json.NewDecoder(reader)
	for {
		var current event
		if err := decoder.Decode(&current); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("decode BEP: %w", err)
		}

		if current.ID.NamedSet != nil && current.NamedSetOfFiles != nil {
			namedSets[current.ID.NamedSet.ID] = *current.NamedSetOfFiles
		}
		if current.Completed == nil {
			continue
		}
		for _, group := range current.Completed.OutputGroups {
			if _, ok := wanted[group.Name]; !ok {
				continue
			}
			if group.Incomplete {
				return nil, fmt.Errorf("output group %q is incomplete", group.Name)
			}
			seenGroups[group.Name] = struct{}{}
			for _, fileSet := range group.FileSets {
				groups[group.Name].roots = append(groups[group.Name].roots, fileSet.ID)
			}
			groups[group.Name].inline = append(groups[group.Name].inline, group.InlineFiles...)
		}
	}

	outputs := make(map[string][]Output, len(groups))
	for name, group := range groups {
		if _, ok := seenGroups[name]; !ok {
			return nil, fmt.Errorf("output group %q was not reported", name)
		}
		files, err := expandNamedSets(namedSets, group.roots)
		if err != nil {
			return nil, fmt.Errorf("read output group %q: %w", name, err)
		}
		files = append(files, outputsFromFiles(group.inline)...)
		outputs[name], err = deduplicate(files)
		if err != nil {
			return nil, fmt.Errorf("read output group %q: %w", name, err)
		}
	}
	return outputs, nil
}

func expandNamedSets(namedSets map[string]namedSetOfFiles, roots []string) ([]Output, error) {
	seenSets := make(map[string]struct{})
	var outputs []Output
	stack := append([]string(nil), roots...)
	for len(stack) > 0 {
		last := len(stack) - 1
		id := stack[last]
		stack = stack[:last]
		if _, ok := seenSets[id]; ok {
			continue
		}
		seenSets[id] = struct{}{}

		set, ok := namedSets[id]
		if !ok {
			return nil, fmt.Errorf("missing named set %q", id)
		}
		outputs = append(outputs, outputsFromFiles(set.Files)...)
		for _, child := range set.FileSets {
			stack = append(stack, child.ID)
		}
	}

	return outputs, nil
}

func outputsFromFiles(files []file) []Output {
	outputs := make([]Output, 0, len(files))
	for _, item := range files {
		outputs = append(outputs, Output{
			Path:              path.Join(append(item.PathPrefix, item.Name)...),
			URI:               item.URI,
			Contents:          item.Contents,
			SymlinkTargetPath: item.SymlinkTargetPath,
			Digest:            item.Digest,
			Length:            item.Length,
		})
	}
	return outputs
}

func deduplicate(files []Output) ([]Output, error) {
	seen := make(map[string]Output, len(files))
	for _, output := range files {
		if previous, ok := seen[output.Path]; ok && previous != output {
			return nil, fmt.Errorf("artifact %q has conflicting metadata", output.Path)
		}
		seen[output.Path] = output
	}

	outputs := make([]Output, 0, len(seen))
	for _, output := range seen {
		outputs = append(outputs, output)
	}
	sort.Slice(outputs, func(i, j int) bool {
		return outputs[i].Path < outputs[j].Path
	})
	return outputs, nil
}
