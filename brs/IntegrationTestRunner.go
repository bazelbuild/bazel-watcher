package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os/exec"
)

var sutBinary string
var testBinary string

func main() {
	flag.StringVar(&sutBinary, "sut_binary", "", "binary of system under test")
	flag.StringVar(&testBinary, "test_binary", "", "test binary to run against system under test")
	flag.Parse()

	port, err := GetEphemeralPort()
	if err != nil {
		panic(err)
	}

	// Bring up the system under test
	sut := exec.Command(sutBinary, "--port", fmt.Sprintf("%d", port), "--nobrowser")
	var sutOut io.ReadCloser
	if sutOut, err = sut.StdoutPipe(); err != nil {
		panic(err)
	}
	go func() {
		if err := sut.Run(); err != nil {
			panic(err)
		}
	}()
	var line string
	if line, err = bufio.NewReader(sutOut).ReadString('\n'); err != nil {
		panic(line)
	}
	fmt.Printf("hello from sut: %s\n", line)
}
