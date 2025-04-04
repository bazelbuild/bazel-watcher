package e2e

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
	"golang.org/x/tools/txtar"
)

var ibazelPath = getiBazelPath()

func getiBazelPath() string {
	path, ok := bazel.FindBinary("cmd/ibazel", "ibazel")
	if !ok {
		panic("Failed to locate binary //ibazel:ibazel, please add it as a data dependency")
	}
	return path
}

func Must(t *testing.T, e error) {
	t.Helper()
	if e != nil {
		t.Fatalf("Error: %s", e)
		// t.Logf("Stack trace:\n%s", string(debug.Stack()))
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

		if strings.HasSuffix(path, ".sh") {
			if err := os.Chmod(path, 0777); err != nil {
				t.Fatalf("Error os.Chmod(%q, 0777): %v", path, err)
			}
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

type Args struct {
	// Main is a text archive containing files in the main workspace.
	// The text archive format is parsed by
	// //go/tools/internal/txtar:go_default_library, which is copied from
	// cmd/go/internal/txtar. If this archive does not contain a WORKSPACE file,
	// a default file will be synthesized.
	Main string

	// SetUp is a function that is executed inside the context of the testing
	// workspace. It is executed once and only once before the beginning of
	// all tests. If SetUp returns a non-nil error, execution is halted and
	// tests cases are not executed.
	SetUp func() error
}

func TestMain(m *testing.M, args Args) {
	print := flag.Bool("print", false, "print out the directory listing before running the test")

	var skipGeneratingWorkspace bool
	ar := txtar.Parse([]byte(args.Main))
	for _, f := range ar.Files {
		switch f.Name {
		case "WORKSPACE":
			skipGeneratingWorkspace = true
		case ".bazelrc":
			fmt.Printf("Do not specify a %q file in your test case. It is not allowed in order to prevent difficult to debug things, like having a .bazelrc file in the test user's home directory be accidentally detected as the one to use in here.", f.Name)
			os.Exit(1)
		}
	}

	var additional string
	if !skipGeneratingWorkspace {
		additional = `
-- WORKSPACE --
# Workspace intentionally left empty
`
	}

	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: args.Main + additional,
		SetUp: func() error {
			if *print {
				if err := filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
					if d.IsDir() {
						return nil
					}
					data, _ := os.ReadFile(path)
					fmt.Fprintf(os.Stderr, "-- %s --\n%s\n", path, string(data))
					return nil
				}); err != nil {
					return fmt.Errorf("WalkDir(.): %w", err)
				}
			}

			if args.SetUp != nil {
				if err := args.SetUp(); err != nil {
					return err
				}
			}

			return nil
		},
	})

	// rules_go does not adhere to the written documentation contract for
	// TestMain. The test main function invokes os.Exit in a deferred block. It
	// can not be nested, it is not reentrant and no code may follow it in this
	// function.
}
