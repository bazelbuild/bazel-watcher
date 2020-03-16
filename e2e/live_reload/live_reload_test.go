package live_reload

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/bazelbuild/rules_webtesting/go/webtest"
	"github.com/google/go-cmp/cmp"
	"github.com/gorilla/websocket"
	"github.com/tebeka/selenium"

	"github.com/bazelbuild/bazel-watcher/e2e"
	"github.com/bazelbuild/bazel-watcher/e2e/live_reload/example_client/data"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

func TestMain(m *testing.M) {
	// mainFiles is a string that describes all the files used in this test in
	// txtar format as described by cmd/go/internal/txtar.
	mainFiles := ""

	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("os.Getwd() error: %v\n", err)
		os.Exit(1)
	}

	exampleClientPath := filepath.Join(wd, "e2e", "live_reload", "example_client")
	if err := filepath.Walk(exampleClientPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Nothing to do here. Dirs are made implicitly in txtar files.
		if info.IsDir() {
			return nil
		}

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("ioutil.ReadFile(%q) err: %v", path, err)
		}

		shortPath := strings.TrimPrefix(path, wd+string(filepath.Separator))

		mainFiles += fmt.Sprintf(`
-- %s --
%s
`, shortPath, content)

		return nil
	}); err != nil {
		fmt.Printf("filepath.Walk() error: %v\n", err)
		os.Exit(1)
	}

	bazel_testing.TestMain(m, bazel_testing.Args{
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
		t.Fatalf("conn.ReadMessage(): %s", err)
	}

	got := strings.TrimSpace(string(v))

	if !reflect.DeepEqual(got, want) {
		t.Errorf("websocket read diff: got [%v], want [%v]", got, want)
	}
}

func getServing(t *testing.T, ibazel *e2e.IBazelTester) string {
	t.Helper()

	ibazel.ExpectOutput("Serving: http://.+:\\d+")
	out := ibazel.GetOutput()

	if out == "" {
		t.Fatal("Output was empty. Expected at least some output")
	}

	servingURL := out[len("Serving: "):]

	// Ensure URL validity by parsing it.
	_, err := url.ParseRequestURI(servingURL)
	if err != nil {
		t.Error(err)
	}

	return servingURL
}

func getConfig(t *testing.T, servingURL string) *data.Config {
	t.Helper()

	configResp, err := http.Get(servingURL + "/config")
	if err != nil {
		t.Fatalf("http.Get(%s/config) error: %v", servingURL, err)
	}
	defer configResp.Body.Close()

	configBytes, err := ioutil.ReadAll(configResp.Body)

	var config data.Config
	if err := json.Unmarshal(configBytes, &config); err != nil {
		t.Errorf("json.Unmarshal(configBytes, config) error: %v", err)
	}
	return &config
}

func expectActions(t *testing.T, servingURL string, cmd ...string) {
	t.Helper()

	cfg := getConfig(t, servingURL)
	if diff := cmp.Diff(cfg.Commands, cmd); diff != "" {
		t.Errorf("expectedAction diff: %v", diff)
	}
}

func TestLiveReload(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//e2e/live_reload/example_client:live_reload")
	defer ibazel.Kill()

	servingURL := getServing(t, ibazel)
	config := getConfig(t, servingURL)

	url, err := url.ParseRequestURI(config.LiveReloadURL)
	if err != nil {
		t.Error(err)
	}

	resp, err := http.Get(config.LiveReloadURL)
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
		t.Errorf("websocket.Dial() error: %v", err)
	}

	// Send the hello message to the client.
	hello := liveReloadHello{
		Command:   "hello",
		Protocols: []string{"http://livereload.com/protocols/official-7"},
	}
	if err = conn.WriteJSON(hello); err != nil {
		t.Errorf("conn.WriteJSON() error: %v", err)
	}

	// Verify the hello message
	verify(t, conn, `{"command":"hello","protocols":["http://livereload.com/protocols/official-7","http://livereload.com/protocols/official-8","http://livereload.com/protocols/official-9","http://livereload.com/protocols/2.x-origin-version-negotiation","http://livereload.com/protocols/2.x-remote-control"],"serverName":"live reload"}`)

	e2e.MustWriteFile(t, "test.txt", "2")
	// TODO: RM this line, I don't think it adds any value
	// DO NOT SUBMIT with this line
	ibazel.ExpectOutput("Serving: http://.+:\\d+")

	verify(t, conn, `{"command":"reload","path":"reload","liveCSS":true}`)

	e2e.MustWriteFile(t, "test.txt", "3")
	// DO NOT SUBMIT with this line
	ibazel.ExpectOutput("Serving: http://.+:\\d+")

	verify(t, conn, `{"command":"reload","path":"reload","liveCSS":true}`)

	t.Logf("Output:\n%v", ibazel.GetOutput())
	t.Logf("Error:\n%v", ibazel.GetError())
}

func TestNoLiveReload(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//e2e/live_reload/example_client:no_live_reload")
	defer ibazel.Kill()

	// Expect there to not be a url since this is the negative test case.
	servingURL := getServing(t, ibazel)
	config := getConfig(t, servingURL)
	if config.LiveReloadURL != "" {
		t.Errorf("Expected LiveReloadURL to be empty, was %q", config.LiveReloadURL)
	}
}

func TestBrowserLiveReload(t *testing.T) {
	if _, ok := os.LookupEnv("WEB_TEST_WEBDRIVER_SERVER"); !ok {
		t.Skipf("WEB_TEST_WEBDRIVER_SERVER was not set and it must be to run this test. To run this please run the :web_test target in the same directory")
	}

	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//e2e/live_reload/example_client:live_reload")
	defer ibazel.Kill()

	// While the build is happening, start up a webbrowser.
	wd, err := webtest.NewWebDriverSession(selenium.Capabilities{})
	if err != nil {
		t.Fatalf("webtest.NewWebDriverSession() err: %v", err)
		t.Fatal(err)
	}

	// Now blockingly wait for the URL
	servingURL := getServing(t, ibazel)

	t.Logf("Serving URL: %v", servingURL)

	if err := wd.Get(servingURL); err != nil {
		t.Fatalf("wd.Get(%q) err: %v", servingURL, err)
	}

	el, err := wd.FindElement(selenium.ByCSSSelector, "body")
	if err != nil {
		t.Fatalf("wd.FindElement(body) error: %v", err)
	}

	if text, err := el.Text(); err != nil {
		t.Fatalf("el.Text() error: %v", err)
	} else if text != "1" {
		t.Fatalf("el.Text() should be 1 but was %q", text)
	}
	e2e.MustWriteFile(t, "test.txt", "2")

	if text, err := el.Text(); err != nil {
		t.Fatalf("el.Text() error: %v", err)
	} else if text != "2" {
		t.Fatalf("el.Text() should be 2 but was %q", text)
	}

	if err := wd.Quit(); err != nil {
		t.Logf("Error quitting webdriver: %v", err)
	}

	e2e.MustWriteFile(t, "test.txt", "3")
}
