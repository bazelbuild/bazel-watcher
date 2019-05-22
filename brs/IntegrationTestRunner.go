package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
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
	portString := fmt.Sprintf("%d", port)

	// Bring up the system under test
	sut := exec.Command(sutBinary, "--port", portString, "--nobrowser")
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

	// Run the test binary
	testDone := make(chan int)
	test := exec.Command(testBinary, "--backend_port", portString)
	go func() {
		if err := test.Run(); err != nil {
			log.Printf("test binary %v exited with %v", testBinary, err)
			testDone <- 1 // TODO propagate actual status (cast to ExitError)
		} else {
			testDone <- 0
		}
	}()

	status := <-testDone
	os.Exit(status)
}
