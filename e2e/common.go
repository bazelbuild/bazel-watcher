package e2e

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

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

var ibazelPath = getiBazelPath()

func getiBazelPath() string {
	suffix := ""
	// Windows expects executables to end in .exe
	if runtime.GOOS == "windows" {
		suffix = ".exe"
	}
	return GetPath(fmt.Sprintf("ibazel/%s_%s_pure_stripped/ibazel%s", runtime.GOOS, runtime.GOARCH, suffix))
}
