package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestEnsureBazel8Compatibility(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		setup    func(tempDir string) error
		expected []string
	}{
		{
			name:     "respects existing bzlmod flag",
			args:     []string{"--enable_bzlmod"},
			setup:    func(tempDir string) error { return nil },
			expected: []string{"--enable_bzlmod"},
		},
		{
			name:     "respects existing workspace flag",
			args:     []string{"--enable_workspace"},
			setup:    func(tempDir string) error { return nil },
			expected: []string{"--enable_workspace"},
		},
		{
			name: "MODULE.bazel uses bzlmod",
			args: []string{},
			setup: func(tempDir string) error {
				return os.WriteFile(filepath.Join(tempDir, "MODULE.bazel"), []byte("module(name = \"test\")"), 0644)
			},
			expected: []string{"--enable_bzlmod"},
		},
		{
			name: "WORKSPACE uses workspace mode",
			args: []string{},
			setup: func(tempDir string) error {
				return os.WriteFile(filepath.Join(tempDir, "WORKSPACE"), []byte("workspace(name = \"test\")"), 0644)
			},
			expected: []string{"--enable_workspace"},
		},
		{
			name: "local_repository forces workspace mode",
			args: []string{},
			setup: func(tempDir string) error {
				content := `local_repository(name = "test", path = "/path")`
				return os.WriteFile(filepath.Join(tempDir, "WORKSPACE"), []byte(content), 0644)
			},
			expected: []string{"--enable_workspace"},
		},
		{
			name:     "defaults to bzlmod",
			args:     []string{},
			setup:    func(tempDir string) error { return nil },
			expected: []string{"--enable_bzlmod"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for test
			tempDir, err := os.MkdirTemp("", "bazel_compat_test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Change to temp directory
			oldDir, err := os.Getwd()
			if err != nil {
				t.Fatalf("Failed to get current dir: %v", err)
			}
			defer os.Chdir(oldDir)

			if err := os.Chdir(tempDir); err != nil {
				t.Fatalf("Failed to change to temp dir: %v", err)
			}

			// Run test setup
			if err := tt.setup(tempDir); err != nil {
				t.Fatalf("Test setup failed: %v", err)
			}

			// Run the function
			result := EnsureBazel8Compatibility(tt.args)

			// Check result
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("EnsureBazel8Compatibility() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestHasLocalRepository(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name:     "detects local_repository",
			content:  `local_repository(name = "test", path = "/path")`,
			expected: true,
		},
		{
			name:     "ignores other rules",
			content:  `http_archive(name = "test", url = "...")`,
			expected: false,
		},
		{
			name:     "handles missing file",
			content:  "", // Will test with no file
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := os.MkdirTemp("", "test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			oldDir, _ := os.Getwd()
			defer os.Chdir(oldDir)
			os.Chdir(tempDir)

			if tt.content != "" {
				os.WriteFile("WORKSPACE", []byte(tt.content), 0644)
			}
			// For empty content, we deliberately don't create the file

			if got := hasLocalRepository(); got != tt.expected {
				t.Errorf("got %v, want %v", got, tt.expected)
			}
		})
	}
}
