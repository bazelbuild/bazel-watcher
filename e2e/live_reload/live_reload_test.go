// Copyright 2017 The Bazel Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package live_reload

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	bazel "github.com/bazelbuild/bazel-integration-testing/go"
	"github.com/bazelbuild/bazel-watcher/e2e"
	"github.com/gorilla/websocket"
)

type liveReloadHello struct {
	Command   string   `json:"command"`
	Protocols []string `json:"protocols"`
}

const printLivereload = `printf "Live reload url: ${IBAZEL_LIVERELOAD_URL}"`

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

func verify(t *testing.T, conn *websocket.Conn, expected interface{}) {
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))

	_, v, err := conn.ReadMessage()
	m := strings.TrimSpace(string(v))
	t.Logf("v: %s, err: %s\n", m, err)
	if err != nil {
		t.Errorf("Error ReadJSONing from websocket: %s", err)
	}

	assertEqual(t, expected, m, "Expected response match")
}

func TestLiveReload(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchFileWithMode("test.sh", printLivereload, 0777))
	must(t, b.ScratchFile("test.txt", "1"))
	must(t, b.ScratchFile("BUILD", `
sh_binary(
	name = "live_reload",
	srcs = ["test.sh"],
	data = ["test.txt"],
	tags = ["ibazel_live_reload"],
)
`))

	ibazel := e2e.NewIBazelTester(t, b)
	ibazel.Run("//:live_reload")
	defer ibazel.Kill()

	ibazel.ExpectOutput("Live reload url: http://.+:\\d+")
	out := ibazel.GetOutput()
	t.Logf("Output: '%s'", out)

	jsUrl := out[len("Live reload url: "):]
	t.Logf("Livereload URL: '%s'", jsUrl)

	url, err := url.ParseRequestURI(jsUrl)
	if err != nil {
		t.Error(err)
	}

	client := &http.Client{}
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

	expectedStart := "(function e(t,n,r)"
	actualStart := bodyString[0:len(expectedStart)]
	if actualStart != expectedStart {
		t.Errorf("Expected js to start with \"%s\" but got \"%s\".", expectedStart, actualStart)
	}

	expectedEnd := "},{}]},{},[8]);"
	actualEnd := bodyString[len(bodyString)-len(expectedEnd):]
	if actualEnd != expectedEnd {
		t.Errorf("Expected js to end with \"%s\" but got \"%s\".", expectedEnd, actualEnd)
	}

	wsUrl := "ws://" + url.Hostname() + ":" + url.Port() + "/livereload"
	t.Logf("wsUrl: %s", wsUrl)
	conn, _, err := websocket.DefaultDialer.Dial(wsUrl, map[string][]string{})
	if err != nil {
		t.Error(err)
	}

	// Send the hello message to the client.
	hello := liveReloadHello{
		Command:   "hello",
		Protocols: []string{"http://livereload.com/protocols/official-7"},
	}
	if err = conn.WriteJSON(hello); err != nil {
		t.Error(err)
	}

	// Verify the hello message
	verify(t, conn, `{"command":"hello","protocols":["http://livereload.com/protocols/official-7","http://livereload.com/protocols/official-8","http://livereload.com/protocols/official-9","http://livereload.com/protocols/2.x-origin-version-negotiation","http://livereload.com/protocols/2.x-remote-control"],"serverName":"live reload"}`)

	must(t, b.ScratchFile("test.txt", "2"))
	ibazel.ExpectOutput("Live reload url: http://.+:\\d+")

	verify(t, conn, `{"command":"reload","path":"reload","liveCSS":true}`)

	must(t, b.ScratchFile("test.txt", "3"))
	ibazel.ExpectOutput("Live reload url: http://.+:\\d+")

	verify(t, conn, `{"command":"reload","path":"reload","liveCSS":true}`)
}

func TestNoLiveReload(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchFileWithMode("test.sh", printLivereload, 0777))
	must(t, b.ScratchFile("BUILD", `
sh_binary(
	name = "no_live_reload",
	srcs = ["test.sh"],
)
`))

	ibazel := e2e.NewIBazelTester(t, b)
	ibazel.Run("//:no_live_reload")
	defer ibazel.Kill()

	// Expect there to not be a url since this is the negative test case.
	ibazel.ExpectOutput("Live reload url: $")
}
