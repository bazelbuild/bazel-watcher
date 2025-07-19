package utils

import (
	"os"
	"strings"
)

// EnsureBazel8Compatibility adds appropriate --enable_bzlmod or --enable_workspace
// flags to bazel arguments when running with Bazel 8+, which requires explicit
// selection between bzlmod and WORKSPACE modes.
func EnsureBazel8Compatibility(args []string) []string {
	// If user already specified a mode, respect their choice
	for _, arg := range args {
		if strings.HasPrefix(arg, "--enable_bzlmod") || strings.HasPrefix(arg, "--enable_workspace") {
			return args
		}
	}

	// Check for local_repository usage (forces WORKSPACE mode)
	if hasLocalRepository() {
		return append(args, "--enable_workspace")
	}

	// Auto-detect based on project files in current directory
	if _, err := os.Stat("MODULE.bazel"); err == nil {
		return append(args, "--enable_bzlmod")
	}
	if _, err := os.Stat("WORKSPACE"); err == nil {
		return append(args, "--enable_workspace")
	}
	if _, err := os.Stat("WORKSPACE.bazel"); err == nil {
		return append(args, "--enable_workspace")
	}

	// Default to bzlmod for new projects
	return append(args, "--enable_bzlmod")
}

func hasLocalRepository() bool {
	workspaceData, err := os.ReadFile("WORKSPACE")
	if err != nil {
		return false // No WORKSPACE file or can't read it
	}

	for _, line := range strings.Split(string(workspaceData), "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "local_repository(") {
			return true
		}
	}
	return false
}
