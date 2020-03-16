package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net"
	"net/http"
	"os"

	"github.com/bazelbuild/rules_go/go/tools/bazel"

	"github.com/bazelbuild/bazel-watcher/e2e/live_reload/example_client/data"
)

var (
	indexTpl = template.Must(template.New("index").Parse(`
<html>
    <head>
        <script src="{{ .LiveReloadURL }}"></script>
    </head>
    <body>
        <p>{{ .Number }}</p>
    </body>
</html>`))
)

func main() {
	cfg := &data.Config{}

	if liveReloadURL, ok := os.LookupEnv("IBAZEL_LIVERELOAD_URL"); ok {
		cfg.LiveReloadURL = liveReloadURL
	}

	go func() {
		scan := bufio.NewScanner(os.Stdin)
		for scan.Scan() {
			cfg.Commands = append(cfg.Commands, scan.Text())
			fmt.Printf("Got command: %v\n", scan.Text())
		}
	}()

	s := &http.Server{}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
		numPath, err := bazel.Runfile("e2e/live_reload/example_client/test.txt")
		if err != nil {
			fmt.Fprintf(resp, "bazel.Runfile(\"e2e/live_reload/example_client/test.txt\") error: %v", err)
		}

		num, err := ioutil.ReadFile(numPath)
		if err != nil {
			fmt.Fprintf(resp, "ioutil.ReadFile(%s) error: %v", numPath, err)
		}

		indexTpl.Execute(resp, map[string]interface{}{
			"LiveReloadURL": cfg.LiveReloadURL,
			"Number":        string(num),
		})
	})
	mux.HandleFunc("/config", func(resp http.ResponseWriter, req *http.Request) {
		e := json.NewEncoder(resp)
		e.SetIndent("", "  ")
		e.Encode(cfg)
	})

	s.Handler = mux

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Printf("net.Listen(\"127.0.0.1:0\") error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Serving: http://%s", listener.Addr().String())
	if err := s.Serve(listener); err != nil {
		fmt.Printf("s.Serve(%v) error: %v\n", listener.Addr().String(), err)
		os.Exit(1)
	}
}
