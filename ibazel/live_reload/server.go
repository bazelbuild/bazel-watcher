// Copyright 2017 The Bazel Authors. All rights reserved.
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

package live_reload

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	blaze_query "github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf"
	"github.com/jaschaephraim/lrserver"
)

var noLiveReload = flag.Bool("nolive_reload", false, "Disable JavaScript live reload support")

type LiveReloadServer struct {
	lrserver *lrserver.Server
}

func New() *LiveReloadServer {
	return &LiveReloadServer{}
}

func (l *LiveReloadServer) Initialize() {}

func (l *LiveReloadServer) startLiveReloadServer() {
	if l.lrserver != nil {
		return
	}

	port := lrserver.DefaultPort
	for ; port < lrserver.DefaultPort+100; port++ {
		if testPort(port) {
			l.lrserver = lrserver.New("live reload", port)
			// Live reload server shouldn't log.
			l.lrserver.SetStatusLog(log.New(os.Stderr, "", 0))
			go func() {
				err := l.lrserver.ListenAndServe()
				if err != nil {
					fmt.Fprintf(os.Stderr, "Live reload server failed to start: %v\n", err)
				}
			}()
			url := fmt.Sprintf("http://localhost:%d/livereload.js?snipver=1", port)
			os.Setenv("IBAZEL_LIVERELOAD_URL", url)
			return
		}
	}
	fmt.Fprintf(os.Stderr, "Could not find open port for live reload server\n")
}
func (l *LiveReloadServer) Cleanup() {
	if l.lrserver != nil {
		l.lrserver.Close()
	}
}

func (l *LiveReloadServer) TargetDecider(rule *blaze_query.Rule) {
	for _, attr := range rule.Attribute {
		if *attr.Name == "tags" && *attr.Type == blaze_query.Attribute_STRING_LIST {
			if contains(attr.StringListValue, "ibazel_live_reload") {
				if *noLiveReload {
					fmt.Fprintf(os.Stderr, "Target requests live_reload but liveReload has been disabled with the -nolive_reload flag.\n")
					return
				}
				l.startLiveReloadServer()
				return
			}
		}
	}
}
func contains(l []string, e string) bool {
	for _, i := range l {
		if i == e {
			return true
		}
	}
	return false
}

func (l *LiveReloadServer) ChangeDetected(changeType string) {}

func (l *LiveReloadServer) BeforeCommand(command string) {}
func (l *LiveReloadServer) AfterCommand(command string, success bool) {
	if l.lrserver != nil {
		fmt.Fprintf(os.Stderr, "Triggering live reload\n")
		l.lrserver.Reload("reload")
	}
}

func testPort(port uint16) bool {
	ln, err := net.Listen("tcp", ":"+strconv.FormatInt(int64(port), 10))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Port %d: %v\n", port, err)
		return false
	}

	ln.Close()
	return true
}
