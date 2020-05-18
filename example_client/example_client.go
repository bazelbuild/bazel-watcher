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
	"strings"
	"time"

	"github.com/bazelbuild/rules_go/go/tools/bazel"

	"github.com/bazelbuild/bazel-watcher/example_client/data"
)

var (
	indexTpl = template.Must(template.New("index").Parse(`
<html>
    <head>
        <script src="{{ .LiveReloadURL }}"></script>
    </head>
    <body>
        {{ .Number }}
    </body>
</html>`))
)

func getData() (string, error) {
	dataPath, err := bazel.Runfile("example_client/data.txt")
	if err != nil {
		return "", fmt.Errorf("bazel.Runfile(\"example_client/data.txt\") error: %v", err)
	}

	num, err := ioutil.ReadFile(dataPath)
	if err != nil {
		return "", fmt.Errorf("ioutil.ReadFile(%s) error: %v", dataPath, err)
	}

	return string(num), nil
}

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
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		num, err := getData()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "getData() error: %v")
			return
		}

		indexTpl.Execute(w, map[string]interface{}{
			"LiveReloadURL": cfg.LiveReloadURL,
			"Number":        string(num),
		})
	})
	mux.HandleFunc("/runfile/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.TrimPrefix(r.URL.Path, "/runfile/")
		dataPath, err := bazel.Runfile(path)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "bazel.Runfile(%q) error: %v", path, err)
			return
		}

		fmt.Fprintf(w, "%s", dataPath)
	})
	mux.HandleFunc("/raw", func(w http.ResponseWriter, r *http.Request) {
		num, err := getData()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "getData() error: %v")
			return
		}

		fmt.Fprintf(w, "%s", num)
	})
	mux.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(cfg)
	})
	mux.HandleFunc("/killkillkill", func(w http.ResponseWriter, r *http.Request) {
		delay := time.Second
		fmt.Fprintf(w, "Killing in %v...", delay)
		go func() {
			time.Sleep(delay)
			os.Exit(1)
		}()
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
