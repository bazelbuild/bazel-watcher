package profiler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

const mainFiles = `
-- BUILD --
sh_binary(
	name = "test",
	srcs = ["test.sh"],
)
-- test.sh --
printf "Profiler url: ${IBAZEL_PROFILER_URL}"
`

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: mainFiles,
	})
}

const printLivereload = `printf "Profiler url: ${IBAZEL_PROFILER_URL}"`

type profileEvent struct {
	Type string `json:"type"`
}

func TestProfiler(t *testing.T) {
	// Make a tempfile the profiler can write to.
	tempFile, err := ioutil.TempFile("", "ibazel_profiler_json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(tempFile.Name()); err != nil {
			t.Logf("os.Remove(%q): %v", tempFile.Name(), err)
		}
	}()

	ibazel := e2e.SetUp(t)
	ibazel.RunWithProfiler("//:test", tempFile.Name())
	defer ibazel.Kill()

	ibazel.ExpectOutput("Profiler url: http://.+:\\d+", 35 * time.Second)
	out := ibazel.GetOutput()
	t.Logf("Output: '%s'", out)

	jsUrl := out[len("Profiler url: "):]
	t.Logf("Profiler URL: '%s'", jsUrl)

	_, err = url.ParseRequestURI(jsUrl)
	if err != nil {
		t.Error(err)
	}

	client := http.Client{}
	resp, err := client.Get(jsUrl)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	bodyString := string(body)

	expectedStart := "// Copyright 2017 The Bazel Authors. All rights reserved."
	actualStart := bodyString[0:len(expectedStart)]
	if actualStart != expectedStart {
		t.Errorf("Expected js to start with \"%s\" but got \"%s\".", expectedStart, actualStart)
	}

	expectedEnd := "})();"
	actualEnd := bodyString[len(bodyString)-len(expectedEnd):]
	if actualEnd != expectedEnd {
		t.Errorf("Expected js to end with \"%s\" but got \"%s\".", expectedEnd, actualEnd)
	}

	profileLog, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	events := compact(strings.Split(string(profileLog), "\n"))
	t.Logf("Profile log: %v", events)

	if len(events) != 3 {
		t.Fatal("Expected 3 events")
	}

	var event profileEvent

	err = json.Unmarshal([]byte(events[0]), &event)
	if err != nil {
		t.Errorf("json.Unmarshal([]byte(%q), &event): %v", events[0], err)
	}
	if event.Type != "IBAZEL_START" {
		t.Errorf("Expected IBAZEL_START, got %q", event.Type)
	}

	err = json.Unmarshal([]byte(events[1]), &event)
	if err != nil {
		t.Errorf("json.Unmarshal([]byte(%q), &event): %v", events[1], err)
	}
	if event.Type != "RUN_START" {
		t.Errorf("Expected RUN_START, got %q", event.Type)
	}

	err = json.Unmarshal([]byte(events[2]), &event)
	if err != nil {
		t.Errorf("json.Unmarshal([]byte(%q), &event): %v", events[2], err)
	}
	if event.Type != "RUN_DONE" {
		t.Errorf("Expected RUN_DONE, got %q", event.Type)
	}
}

func TestNoProfiler(t *testing.T) {
	ibazel := e2e.SetUp(t)
	// Note that there is nothing special about the test that makes it a profile
	// run vs a non-profiling run, only command line arguments to ibazel.
	ibazel.Run([]string{}, "//:test")
	defer ibazel.Kill()

	ibazel.ExpectOutput("Profiler url: $")
}

// compact provided slice to only contain non-empty strings.
func compact(a []string) []string {
	var r []string
	for _, str := range a {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
