package notify_changes

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/e2e"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

const mainFiles = `
-- BUILD --
load("@io_bazel_rules_go//go:def.bzl", "go_binary")
go_binary(
  name = "notify_changes",
	srcs = ["devserver.go"],
  data = ["main.js"],
  tags = ["ibazel_notify_changes"],
)
-- main.js --
1
-- .bazelrc --
build --enable_runfiles
run --enable_runfiles
test --enable_runfiles
-- devserver.go --
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		http.ServeFile(res, req, req.URL.Path[1:])
	})

	ts := httptest.NewServer(mux)
	defer ts.Close()

	fmt.Print(ts.URL)
	select{}
}
`

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: mainFiles,
	})
}

func getScriptContent(t *testing.T, basePath string) string {
	client := &http.Client{}
	resp, err := client.Get(basePath + "/main.js")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	return string(body)
}

func TestNotifyChanges(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//:notify_changes")
	defer ibazel.Kill()
	ibazel.ExpectOutput("http://.+:\\d+", time.Minute)
	out := ibazel.GetOutput()

	if out == "" {
		t.Fatal("Output was empty. Expected at least some output")
	}

	content := getScriptContent(t, out)

	if strings.TrimSpace(content) != "1" {
		t.Fatal("The served file content is wrong:", content)
	}

	e2e.MustWriteFile(t, "main.js", "2")

	content = getScriptContent(t, out)

	if strings.TrimSpace(content) != "2" {
		t.Fatal("The served file content is wrong:", content)
	}
}
