// Copyright 2019 The Bazel Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

var sutBinary string
var testBinary string

func main() {
	flag.StringVar(&sutBinary, "sut_binary", "", "binary of system under test")
	flag.StringVar(&testBinary, "test_binary", "", "test binary to run against system under test")
	flag.Parse()

	port, err := getEphemeralPort()
	if err != nil {
		panic(err)
	}
	portString := fmt.Sprintf("%d", port)

	// Bring up the system under test
	sut := exec.Command(sutBinary, "--port", portString)
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

func getEphemeralPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}
	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
