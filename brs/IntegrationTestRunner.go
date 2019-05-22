package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os/exec"
)

var sutBinary string
var testBinary string

func main() {
	flag.StringVar(&sutBinary, "sut_binary", "", "binary of system under test")
	flag.StringVar(&testBinary, "test_binary", "", "test binary to run against system under test")
	flag.Parse()

	// Bring up the system under test
	sut := exec.Command(sutBinary, "--nobrowser")
	var sutOut io.ReadCloser
	var err error
	if sutOut, err = sut.StdoutPipe(); err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := sut.Run(); err != nil {
			log.Fatal(err)
		}
	}()
	var line string
	if line, err = bufio.NewReader(sutOut).ReadString('\n'); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("hello from sut: %s\n", line)
}
