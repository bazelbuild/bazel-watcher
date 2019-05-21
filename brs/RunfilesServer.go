package main

import (
	"fmt"
	"log"
	"net"
)
import "flag"
import "github.com/pkg/browser"
import "net/http"

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
	// Print a line to stdout. IntegrationTestRunner uses this for synchronization (it won't run the
	// test binary until the system under test prints a line to stdout). For other uses, this is
	// harmless.
	fmt.Printf("listening on %d\n", port)
	if shouldOpenBrowser() {
		browser.OpenURL(fmt.Sprintf("http://localhost:%d/%s", port, index))
	}
	http.Serve(listener, http.FileServer(http.Dir(""))) // serve from runfiles root
}

func shouldOpenBrowser() bool {
	return nobrowser == false
}
