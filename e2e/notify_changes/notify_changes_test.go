package notify_changes

import (
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/e2e/example_client"
)

func TestMain(m *testing.M) {
	example_client.TestMain(m)
}

func TestRestartProcess(t *testing.T) {
	client := example_client.StartLiveReload(t)
	defer client.Cleanup()
	client.Kill(t)

	client.SetData(t, "newdata")

	// Wait for ibazel to successfully restart process
	client.DetectServerParameters(t, time.Second*5)

	stderr := client.IBazel.GetIBazelError()
	if strings.Contains(stderr, "Error writing") ||
		strings.Contains(stderr, "broken pipe") {
		t.Errorf(
			"ibazel tried to write to stdout/stderr of terminated process.\nibazel error output: \n%v",
			stderr)
	}
}
