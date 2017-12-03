package e2e

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"runtime/debug"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

func must(t *testing.T, e error) {
	if e != nil {
		t.Errorf("Error: %s", e)
		debug.PrintStack()
	}
}

func assertNotEqual(t *testing.T, want, got interface{}, msg string) {
	if reflect.DeepEqual(want, got) {
		t.Errorf("Wanted %s, got %s. %s", want, got, msg)
		debug.PrintStack()
	}
}
func assertEqual(t *testing.T, want, got interface{}, msg string) {
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Wanted [%v], got [%v]. %s", want, got, msg)
		debug.PrintStack()
	}
}

func getPath(p string) string {
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
	ibazelPath = getPath(fmt.Sprintf("ibazel/%s_%s_pure_stripped/ibazel", runtime.GOOS, runtime.GOARCH))
}
