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
	PublishConfig map[string]string `json:"publishConfig"`
}

// Note that Version must be assigned to a constant value to allow stamping to
// override.
// See:
//   * https://github.com/bazelbuild/rules_go/issues/2379
//   * https://github.com/golang/go/issues/37369
var Version string = "VERSION_NOT_SET"

func main() {
	if "VERSION_NOT_SET" == Version {
		panic("The version string was not overriden. Please rebuild with --stamp")
	}

	pkg := jsPackage{
		Name:        "@bazel/ibazel",
		Version:     strings.TrimPrefix(strings.TrimPrefix(Version, "v"), "V"),
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
