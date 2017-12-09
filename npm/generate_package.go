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
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
)

type jsPackage struct {
	Name          string            `json:"name"`
	Version       string            `json:"version"`
	Description   string            `json:"description"`
	Bin           map[string]string `json:"bin"`
	Contributors  []string          `json:"contributors"`
	License       string            `json:"license"`
	PublishConfig map[string]string `json:"publish_config"`
}

var unsetVersion = "VERSION_NOT_SET"
var Version string = unsetVersion

func main() {
	if unsetVersion == Version {
		panic("The version string was not overriden. Please rebuild with --stamp")
	}

	pkg := jsPackage{
		Name:        "@bazel/ibazel",
		Version:     Version,
		Description: "node distribution of ibazel",
		Bin: map[string]string{
			"ibazel": "index.js",
		},
		Contributors: []string{},
		License:      "Apache-2.0",
		PublishConfig: map[string]string{
			"access": "public",
		},
	}

	file, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	for _, line := range strings.Split(string(file), "\n") {
		if line != "" && line[0] != '#' {
			pkg.Contributors = append(pkg.Contributors, line)
		}
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "\t")
	err = enc.Encode(pkg)
	if err != nil {
		panic(err)
	}
}
