package selenium

import (
	"reflect"
	"runtime/debug"
	"strings"
	"testing"
	"time"

	bazel "github.com/bazelbuild/bazel-integration-testing/go"
	"github.com/bazelbuild/bazel-watcher/e2e"

	"github.com/gorilla/websocket"

	"github.com/bazelbuild/rules_webtesting/go/webtest"
	"github.com/tebeka/selenium"
)

type liveReloadHello struct {
	Command   string   `json:"command"`
	Protocols []string `json:"protocols"`
}

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
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

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
	must(t, b.ScratchFileWithMode("live_reload.py", string(python), 0777))
	must(t, b.ScratchFile("test.txt", "1"))
	must(t, b.ScratchFile("BUILD", `
py_binary(
	name = "live_reload",
	srcs = ["live_reload.py"],
	data = ["test.txt"],
	tags = ["ibazel_notify_changes", "ibazel_live_reload"],
)
`))

	ibazel := e2e.NewIBazelTester(t, b)
	ibazel.Run("//:live_reload")
	defer ibazel.Kill()

	ibazel.ExpectOutput("Webserver url: http://.+:\\d+")
	out := ibazel.GetOutput()
	t.Logf("Output: '%s'", out)

	url := out[len("Webserver url: "):]
	t.Logf("Webserver URL: '%s'", url)

	wd, err := webtest.NewWebDriverSession(selenium.Capabilities{
		"webdriver.logging.profiler.enabled": true,
		"extendedDebugging":                  true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer wd.Quit()

	if err := wd.Get(url); err != nil {
		t.Error(err)
	}

	elem, err := wd.FindElement(selenium.ByTagName, "body")
	if err != nil {
		t.Error(err)
	}

	text, err := elem.Text()
	if err != nil {
		t.Error(err)
	}
	assertEqual(t, "1", text, "Body text was different")

	must(t, b.ScratchFile("test.txt", "2"))
	time.Sleep(5 * time.Second)
	elem, err = wd.FindElement(selenium.ByTagName, "body")
	if err != nil {
		t.Error(err)
	}

	text, err = elem.Text()
	if err != nil {
		t.Error(err)
	}
	assertEqual(t, "2", text, "Body text was different")

	must(t, b.ScratchFile("test.txt", "3"))
	time.Sleep(5 * time.Second)
	elem, err = wd.FindElement(selenium.ByTagName, "body")
	if err != nil {
		t.Error(err)
	}

	text, err = elem.Text()
	if err != nil {
		t.Error(err)
	}
	assertEqual(t, "3", text, "Body text was different")
}

func TestNoLiveReload(t *testing.T) {
	t.Skip()
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchFileWithMode("no_live_reload.py", string(python), 0777))
	must(t, b.ScratchFile("BUILD", `
py_binary(
	name = "no_live_reload",
	srcs = ["no_live_reload.py"],
)
`))

	ibazel := e2e.NewIBazelTester(t, b)
	ibazel.Run("//:no_live_reload")
	defer ibazel.Kill()

	// Expect there to not be a url since this is the negative test case.
	ibazel.ExpectOutput("Webserver url: $")
}
