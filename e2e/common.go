package e2e

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

var _ = runtime.GOOS

func GetPath(p string) string {
	path, err := bazel.Runfile(p)
	if err != nil {
		panic(err)
	}

	path, err = filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	return path
}

var ibazelPath string

func init() {
	ibazelPath = GetPath(fmt.Sprintf("ibazel/%s_%s_pure_stripped/ibazel", runtime.GOOS, runtime.GOARCH))
}
