package example_client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/bazelbuild/rules_go/go/tools/bazel"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"

	"github.com/bazelbuild/bazel-watcher/example_client/data"
	"github.com/bazelbuild/bazel-watcher/internal/e2e"
)

func TestMain(m *testing.M, extraTxtar ...string) {
	// mainFiles is a string that describes all the files used in this test in
	// txtar format as described by cmd/go/internal/txtar.
	mainFiles := ""

	wd, err := os.Getwd()
	if err != nil {
		fmt.Printf("os.Getwd() error: %v\n", err)
		os.Exit(1)
	}

	exampleClientPath := filepath.Join(wd, "example_client")
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

	for _, f := range extraTxtar {
		mainFiles += f
	}

	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: mainFiles,
	})
}

type ExampleClient struct {
	ibazel   *e2e.IBazelTester
	basePath string
	// dataPath points to the file that needs to be updated in order to have the
	// example client render different data. This is tracked in the object
	// because Windows paths are complex and doing that computation many times
	// would be annoying.
	dataPath string
}

func StartLiveReload(t *testing.T) (client *ExampleClient) {
	t.Helper()

	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//example_client:live_reload")

	dataPath, err := bazel.Runfile("example_client/data.txt")
	if err != nil {
		t.Fatalf("bazel.Runfile(\"example_client/data.txt\") error: %v", err)
	}

	c := &ExampleClient{
		ibazel:   ibazel,
		dataPath: dataPath,
	}
	c.DetectServerParameters(t)

	return c
}

func (c *ExampleClient) fatalf(t *testing.T, msg string, args ...interface{}) {
	t.Helper()

	t.Logf("Out: %v", c.ibazel.GetOutput())
	t.Logf("Error: %v", c.ibazel.GetError())
	t.Logf("iBazel Error: %v", c.ibazel.GetIBazelError())

	t.Fatalf(msg, args...)
}

func (c *ExampleClient) Cleanup() {
	c.ibazel.Kill()
}

func (c *ExampleClient) DetectServerParameters(t *testing.T) {
	c.ibazel.ExpectOutput("Serving: http://.+:\\d+")
	out := c.ibazel.GetOutput()

	if out == "" {
		t.Fatal("Output was empty. Expected at least some output")
	}

	urlCaptureGroup := 1
	r := regexp.MustCompile(`Serving: (?P<url>http://[0-9.]+:[0-9]+)`)
	results := r.FindAllStringSubmatch(out, -1)
	if len(results) == 0 {
		c.fatalf(t, "Expected output to decribe where its serving from. Found nothing.")
	}

	c.basePath = results[len(results)-1][urlCaptureGroup]

	// Ensure URL validity by parsing it.
	_, err := url.ParseRequestURI(c.basePath)
	if err != nil {
		t.Error("url.ParseRequestURI(%q)", c.basePath, err)
	}
}

func (c *ExampleClient) get(t *testing.T, path string) ([]byte, error) {
	t.Helper()

	r, err := http.Get(c.basePath + path)
	if err != nil {
		c.fatalf(t, "http.Get(%s/config) error: %v", c.basePath, err)
	}
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll() error: %v", err)
	}
	return b, nil
}

func (c *ExampleClient) GetConfig(t *testing.T) *data.Config {
	configBytes, err := c.get(t, "/config")
	if err != nil {
		t.Errorf("c.get(\"/config\") error: %v", err)
	}

	var config data.Config
	if err := json.Unmarshal(configBytes, &config); err != nil {
		t.Errorf("json.Unmarshal(configBytes, config) error: %v", err)
	}
	return &config
}

func (c *ExampleClient) Kill(t *testing.T) {
	t.Helper()

	_, err := c.get(t, "/killkillkill")
	if err != nil {
		t.Errorf("c.get(\"/killkillkill\") error: %v", err)
	}
}

func (c *ExampleClient) GetRaw(t *testing.T) string {
	t.Helper()

	rawBytes, err := c.get(t, "/raw")
	if err != nil {
		t.Errorf("c.get(\"/raw\") error: %v", err)
	}

	// Files are often written with newlines at the end and therefore will be
	// returned with extra newlines. To simplify testing against values, strip
	// all whitespace.
	return strings.TrimSpace(string(rawBytes))
}

func (c *ExampleClient) SetData(t *testing.T, data string) {
	t.Helper()
	e2e.MustWriteFile(t, c.dataPath, data)
}
