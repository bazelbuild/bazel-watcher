package profiler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"runtime/debug"
	"strings"
	"testing"

	bazel "github.com/bazelbuild/bazel-integration-testing/go"
	"github.com/bazelbuild/bazel-watcher/e2e"
)

const printLivereload = `printf "Profiler url: ${IBAZEL_PROFILER_URL}"`

type profileEvent struct {
	Type string `json:"type"`
}

func must(t *testing.T, e error) {
	if e != nil {
		t.Errorf("Error: %s", e)
		debug.PrintStack()
	}
}

func assertEqual(t *testing.T, want, got interface{}, msg string) {
	if !reflect.DeepEqual(want, got) {
		t.Errorf("Wanted [%v], got [%v]. %s", want, got, msg)
		debug.PrintStack()
	}
}

func TestProfiler(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchFileWithMode("test.sh", printLivereload, 0777))
	must(t, b.ScratchFile("BUILD", `
sh_binary(
	name = "profiler",
	srcs = ["test.sh"],
)
`))

	tempFile, err := ioutil.TempFile("", "ibazel_profiler_json")
	if err != nil {
		t.Fatal(err)
	}

	ibazel := e2e.NewIBazelTester(t, b)
	ibazel.RunWithProfiler("//:profiler", tempFile.Name())
	defer ibazel.Kill()

	ibazel.ExpectOutput("Profiler url: http://.+:\\d+")
	out := ibazel.GetOutput()
	t.Logf("Output: '%s'", out)

	jsUrl := out[len("Profiler url: "):]
	t.Logf("Profiler URL: '%s'", jsUrl)

	_, err = url.ParseRequestURI(jsUrl)
	if err != nil {
		t.Errorf("Failed to parse: %v", err)
	}

	client := http.Client{}
	resp, err := client.Get(jsUrl)
	if err != nil {
		t.Errorf("Failed to get: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Errorf("Failed to ReadAll: %v", err)
	}

	bodyString := strings.Trim(string(body), "\t \n\r")

	expectedStart := "// Copyright 2017 The Bazel Authors. All rights reserved."
	if !strings.HasPrefix(bodyString, expectedStart) {
		t.Errorf("Expected js to start with \"%s\" but got \"%s\".", expectedStart, bodyString)
	}

	expectedEnd := "})();"
	if !strings.HasSuffix(bodyString, expectedEnd) {
		t.Errorf("Expected js to end with \"%s\" but got \"%s\".", expectedEnd, bodyString)
	}

	profileLog, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		t.Errorf("Failed to ReadFile: %v", err)
	}
	events := compact(strings.Split(string(profileLog), "\n"))
	t.Logf("Profile log: %v", events)

	if len(events) != 3 {
		t.Fatal("Expected 3 events")
	}

	var event profileEvent

	err = json.Unmarshal([]byte(events[0]), &event)
	if err != nil {
		t.Errorf("Failed to decode IBAZEL_START event: %v", err)
	}
	assertEqual(t, event.Type, "IBAZEL_START", "Expected IBAZEL_START")

	err = json.Unmarshal([]byte(events[1]), &event)
	if err != nil {
		t.Errorf("Failed to decode RUN_START event: %v", err)
	}
	assertEqual(t, event.Type, "RUN_START", "Expected RUN_START")

	err = json.Unmarshal([]byte(events[2]), &event)
	if err != nil {
		t.Errorf("Failed to decode RUN_DONE event: %v", err)
	}
	assertEqual(t, event.Type, "RUN_DONE", "Expected RUN_DONE")
}

func TestNoProfiler(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchFileWithMode("test.sh", printLivereload, 0777))
	must(t, b.ScratchFile("BUILD", `
sh_binary(
	name = "no_profiler",
	srcs = ["test.sh"],
)
`))

	ibazel := e2e.NewIBazelTester(t, b)
	ibazel.Run("//:no_profiler")
	defer ibazel.Kill()

	ibazel.ExpectOutput("Profiler url:$")
}

func compact(a []string) []string {
	var r []string
	for _, str := range a {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}
