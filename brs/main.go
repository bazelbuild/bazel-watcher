package main

import (
	"fmt"
	"log"
	"net"
	"os"
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
	// TODO: synchronization
	os.Stdout.Write([]byte(fmt.Sprintf("listening on %v\n", listener.Addr())))
	http.Serve(listener, http.HandlerFunc(hello))
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi!"))
}
