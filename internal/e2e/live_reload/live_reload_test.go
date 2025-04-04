package live_reload

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
	"github.com/gorilla/websocket"
)

const mainFiles = `
-- BUILD --
sh_binary(
  name = "live_reload",
  srcs = ["test.sh"],
  # Add a simple data dependency that you can modify.
  data = ["test.txt"],
  tags = ["ibazel_live_reload"],
)
sh_binary(
	name = "no_live_reload",
	srcs = ["test.sh"],
)
-- test.txt --
1
-- test.sh --
printf "Live reload url: ${IBAZEL_LIVERELOAD_URL}"
`

func TestMain(m *testing.M) {
	e2e.TestMain(m, e2e.Args{
		Main: mainFiles,
	})
}

type liveReloadHello struct {
	Command   string   `json:"command"`
	Protocols []string `json:"protocols"`
}

func assertEqual(t *testing.T, want, got interface{}, msg string) {
	t.Helper()
}

func verify(t *testing.T, conn *websocket.Conn, want interface{}) {
	t.Helper()

	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	_, v, err := conn.ReadMessage()
	if err != nil {
		t.Errorf("conn.ReadMessage(): %s", err)
	}

	got := strings.TrimSpace(string(v))

	if !reflect.DeepEqual(got, want) {
		t.Errorf("websocket read diff: got [%v], want [%v]", got, want)
	}
}

func TestLiveReload(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//:live_reload")
	defer ibazel.Kill()

	ibazel.ExpectOutput("Live reload url: http://.+:\\d+", 35 * time.Second)
	out := ibazel.GetOutput()
	t.Logf("Output: '%s'", out)

	if out == "" {
		t.Fatal("Output was empty. Expected at least some output")
	}

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

	e2e.MustWriteFile(t, "test.txt", "2")
	ibazel.ExpectOutput("Live reload url: http://.+:\\d+")

	verify(t, conn, `{"command":"reload","path":"reload","liveCSS":true}`)

	e2e.MustWriteFile(t, "test.txt", "3")
	ibazel.ExpectOutput("Live reload url: http://.+:\\d+")

	verify(t, conn, `{"command":"reload","path":"reload","liveCSS":true}`)
}

func TestNoLiveReload(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//:no_live_reload")
	defer ibazel.Kill()

	// Expect there to not be a url since this is the negative test case.
	ibazel.ExpectOutput("Live reload url: $")
}
