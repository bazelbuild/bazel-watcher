package main

import (
	"fmt"
	"log"
	"net"
)
import "flag"
import "net/http"

func main() {
	port := flag.Int("port", 0, "port to listen on. if not given, an ephemeral port will be chosen")
	//index := flag.String("index", "", "page to visit in the system's default browser when the server is up. If not given, the browser will not be launched")
	flag.Parse()

	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatal(err)
	}
	// Print a line to stdout. IntegrationTestRunner uses this for synchronization (it won't run the
	// test binary until the system under test prints a line to stdout). For other uses, this is
	// harmless.
	fmt.Println("listening on", listener.Addr())
	http.Serve(listener, http.FileServer(http.Dir(""))) // serve from runfiles root
}