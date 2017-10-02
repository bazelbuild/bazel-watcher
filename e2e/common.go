package e2e

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"runtime/debug"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

func notError(t *testing.T, e error) {
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
	return path
}

var bazelPath string
var ibazelPath string

const mainGoPath = "e2e/simple/main.go"
const BUILDPath = "e2e/simple/BUILD.bazel"

func init() {
	bazelPath = getPath("e2e/bazel/bazel")
	ibazelPath = getPath("ibazel/ibazel")

	// Create the files that are actually watched by the test.
	manipulateSourceFile(0)
	manipulateBUILDFile(0)
}

func manipulateSourceFile(seed int) {
	err := ioutil.WriteFile(mainGoPath, []byte(fmt.Sprintf("Not go code %v", seed)), 0755)
	if err != nil {
		panic(err)
	}
}
func manipulateBUILDFile(seed int) {
	err := ioutil.WriteFile(BUILDPath, []byte(fmt.Sprintf("Not BUILD code %v", seed)), 0755)
	if err != nil {
		panic(err)
	}
}
