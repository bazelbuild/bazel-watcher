package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-watcher/internal/ibazel/log"
	"github.com/google/go-cmp/cmp"
)

func TestAppleCaseInsensitivity(t *testing.T) {
	log.SetTesting(t)

	tests := map[string]struct {
		// startingWD to start the test in when evaluating FindWorkspace
		startingWD string
		// wantPath should be the relative path the workspace is found in.
		wantPath string
		// dirs are recursively created to allow both dirs and files to exist.
		dirs []string
		// path to create the workspace in. Must be created by dirs.
		workspacePath string
		// files are touched and left empty.
		files []string
		err   bool
	}{
		"simple": {
			startingWD:    "",
			wantPath:      "",
			dirs:          []string{},
			workspacePath: "/WORKSPACE",
			files:         []string{},
			err:           false,
		},
		"simple with bazel extension": {
			startingWD:    "",
			wantPath:      "",
			dirs:          []string{},
			workspacePath: "/WORKSPACE.bazel",
			files:         []string{},
			err:           false,
		},
		"simple with bzlmod extension": {
			startingWD:    "",
			wantPath:      "",
			dirs:          []string{},
			workspacePath: "/WORKSPACE.bzlmod",
			files:         []string{},
			err:           false,
		},
		"simple with MODULE.bazel": {
			startingWD:    "",
			wantPath:      "",
			dirs:          []string{},
			workspacePath: "/MODULE.bazel",
			files:         []string{},
			err:           false,
		},
		"no workspace": {
			startingWD: "c/d",
			wantPath:   "",
			dirs: []string{
				"a/b",
				"c/d",
			},
			workspacePath: "/a/WORKSPACE",
			files:         []string{},
			err:           true,
		},
		"no workspace in workspace-named path": {
			startingWD: "c/WORKSPACE",
			wantPath:   "/a",
			dirs: []string{
				"a/b",
				"c/WORKSPACE",
			},
			workspacePath: "/a/WORKSPACE",
			files:         []string{},
			err:           true,
		},
		"no workspace in MODULE.bazel-named path": {
			startingWD: "c/MODULE.bazel",
			wantPath:   "",
			dirs: []string{
				"a/b",
				"c/MODULE.bazel",
			},
			workspacePath: "/a/WORKSPACE",
			files:         []string{},
			err:           true,
		},
		"workspace nested in workspace-named path": {
			// this is intended to catch case-insensitive Macs but mimics the
			// case-insensitive match with a case-sensitive match of the dir
			// the "WORKSPACE" dirname should not early-quit the search, allowing
			// us to find /c/MODULE.bazel
			startingWD: "c/d/WORKSPACE",
			wantPath:   "/c",
			dirs: []string{
				"a/b",
				"c/d/WORKSPACE",
			},
			workspacePath: "/c/MODULE.bazel",
			files:         []string{},
			err:           false,
		},
		"nested workspace": {
			startingWD: "a/workspace",
			wantPath:   "",
			dirs: []string{
				"a/workspace",
			},
			workspacePath: "/WORKSPACE",
			files:         []string{},
			err:           false,
		},
		"nested workspace in directory containing nested workspace": {
			startingWD: "a",
			wantPath:   "",
			dirs: []string{
				"a/workspace",
			},
			workspacePath: "/WORKSPACE",
			files:         []string{},
			err:           false,
		},
	}

	startDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("pwd: %v", err)
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			defer os.Chdir(startDir)
			base := t.TempDir()

			// t.TempDir may return a path that isn't the place you end up when you
			// cd to it, so cd there and get that value.
			if err := os.Chdir(base); err != nil {
				t.Fatalf("cd %q", base)
			}
			var err error
			if base, err = os.Getwd(); err != nil {
				t.Fatalf("pwd: %v", err)
			}

			for _, dir := range test.dirs {
				path := filepath.Join(base, dir)
				if err := os.MkdirAll(path, 0755); err != nil {
					t.Fatalf("MkdirAll(%q): %v", path, err)
				}
			}

			baseWorkspace := filepath.Join(base, test.workspacePath)
			if err := os.WriteFile(baseWorkspace, []byte{}, 0644); err != nil {
				t.Fatalf("os.WriteFile(%q): %v", baseWorkspace, err)
			}

			for _, file := range test.files {
				path := filepath.Join(base, file)
				dir := filepath.Dir(path)
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("MkdirAll(%q): %v", dir, err)
				}
				if err := os.WriteFile(path, []byte{}, 0644); err != nil {
					t.Fatalf("os.WriteFile(%q): %v", path, err)
				}
			}

			startingWD := filepath.Join(base, test.startingWD)
			if err := os.Chdir(startingWD); err != nil {
				t.Fatalf("cd %q", startingWD)
			}

			wf := &MainWorkspace{}
			got, err := wf.FindWorkspace()
			if test.err {
				if err == nil {
					t.Errorf("FindWorkspace() got nil want an error")
				}
			} else {
				if err != nil {
					t.Fatalf("FindWorkspace(): %v", err)
				}

				got = strings.TrimPrefix(got, base)
				if diff := cmp.Diff(got, test.wantPath); diff != "" {
					t.Errorf("FindWorkspace() diff (-got,+want):\n%s", diff)
				}
			}
		})
	}
}
