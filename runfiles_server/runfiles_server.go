package main

import (
	"flag"
	"fmt"
	"github.com/pkg/browser"
	"log"
	"net"
	"net/http"
	"os"
)

var port int
var nobrowser bool
var index string

func main() {
	flag.IntVar(&port, "port", 0, "port to listen on. if not given, an ephemeral port will be chosen")
	flag.BoolVar(&nobrowser, "nobrowser", false, `Disables opening the browser. The default behavior
of RunfilesServer is to open a browser to the page given by --index. Pass --nobrowser if this
behavior is not appropriate.`)
	flag.StringVar(&index, "index", "", `page to visit in the system's default browser when the
server is up. If not given, the browser will not be launched`)
	flag.Parse()
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatal(err)
	}
	// Find the actual port in case the passed-in value was ephemeral
	port = listener.Addr().(*net.TCPAddr).Port
	handler := &liveReloadSnippetInjectingHandler{
		Handler: http.FileServer(http.Dir(".")),
		snippet: maybeFormatLiveReloadSnippet(),
	}
	// Print a line to stdout. IntegrationTestRunner uses this for synchronization (it won't run the
	// test binary until the system under test prints a line to stdout). For other uses, this is
	// harmless.
	fmt.Printf("listening on %d\n", port)
	if shouldOpenBrowser() {
		go browser.OpenURL(fmt.Sprintf("http://localhost:%d/%s", port, index))
	}
	http.Serve(listener, handler)
}

func maybeFormatLiveReloadSnippet() []byte {
	// If the this server is being run under ibazel as part of a target that has the tag
	// `ibazel_live_reload`, ibazel will set the IBAZEL_LIVERELOAD_URL environment variable.
	liveReloadUrl := os.Getenv("IBAZEL_LIVERELOAD_URL")
	if len(liveReloadUrl) > 0 {
		return []byte(fmt.Sprintf("<script src=\"%s\"></script>", liveReloadUrl))
	}
	return nil
}

// When the binary is invoked with --index foo.html, we usually want to open a browser to foo.html.
// But also allow for an explicit --nobrowser override, so that it is possible to bazel run a
// serve() target without launching a browser and without changing the target's attributes:
// `bazel run :some_serve_target -- --nobrowser`.
func shouldOpenBrowser() bool {
	return index != "" && nobrowser == false
}
