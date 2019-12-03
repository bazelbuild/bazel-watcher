package e2e

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"testing"

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

func Must(t *testing.T, e error) {
	t.Helper()
	if e != nil {
		t.Fatalf("Error: %s", e)
		//t.Logf("Stack trace:\n%s", string(debug.Stack()))
	}
}

func MustMkdir(t *testing.T, path string, mode ...os.FileMode) {
	t.Helper()

	m := os.FileMode(0777)
	if len(mode) == 1 {
		m = mode[0]
	}

	if err := os.MkdirAll(path, m); err != nil {
		t.Fatalf("os.MkdirAll(%q, %d): %v", path, m, err)
	}
}

func MustWriteFile(t *testing.T, path, content string, mode ...os.FileMode) {
	t.Helper()

	m := os.FileMode(0777)
	if len(mode) == 1 {
		m = mode[0]
	}

	if err := ioutil.WriteFile(path, []byte(content), m); err != nil {
		t.Fatalf("ioutil.WriteFile(%q, []byte(%q), 0777): %v", path, content, err)
	}
}

func SetExecuteBit(t *testing.T) {
	// Before doing anything, set the executable bit on all the .sh files
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if err := os.Chmod(path, 0777); err != nil {
			t.Fatalf("Error os.Chmod(%q, 0777): %v", path, err)
		}

		return nil
	})

	if err != nil {
		t.Fatalf("filpath.Walk(): %v", err)
	}
}

func SetUp(t *testing.T) *IBazelTester {
	SetExecuteBit(t)
	return NewIBazelTester(t)
}
