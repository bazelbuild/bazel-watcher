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

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bazelbuild/bazel-watcher/bazel"
	"github.com/fsnotify/fsnotify"
)

var osExit = os.Exit
var bazelNew = bazel.New
var execCommand = exec.Command

type State string

const (
	DEBOUNCE_QUERY State = "DEBOUNCE_QUERY"
	QUERY          State = "QUERY"
	WAIT           State = "WAIT"
	DEBOUNCE_RUN   State = "DEBOUNCE_RUN"
	RUN            State = "RUN"
	QUIT           State = "QUIT"
)

const debounceDuration = 100 * time.Millisecond
const sourceQuery = "kind('source file', deps(set(%s)))"
const buildQuery = "buildfiles(deps(set(%s)))"

type IBazel struct {
	b *bazel.Bazel

	cmd       *exec.Cmd
	args      []string
	bazelArgs []string

	sigs           chan os.Signal // Signals channel for the current process
	interruptCount int

	buildFileWatcher  *fsnotify.Watcher
	sourceFileWatcher *fsnotify.Watcher

	sourceEventHandler *SourceEventHandler

	state State
}

func New() (*IBazel, error) {
	i := &IBazel{}
	err := i.setup()
	if err != nil {
		return nil, err
	}

	i.sigs = make(chan os.Signal, 1)
	signal.Notify(i.sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			i.handleSignals()
		}
	}()

	return i, nil
}

func (i *IBazel) handleSignals() {
	// Got an OS signal (SIGINT, SIGKILL).
	sig := <-i.sigs
	fmt.Printf("XXXXX: GOT A SIGNAL: %s\n", sig)

	switch sig {
	case syscall.SIGINT:
		if i.subprocessRunning() {
			fmt.Printf("XXXXX: PROCESS IS RUNNING %s\n", sig)
			fmt.Printf("\nSubprocess killed from getting SIGINT\n")
			i.kill()
		} else {
			fmt.Printf("XXXXX: PROCESS IS NOT RUNNING %s\n", sig)
			osExit(3)
		}
		break
	case syscall.SIGKILL:
		if i.subprocessRunning() {
			fmt.Printf("\nSubprocess killed from getting SIGKILL\n")
			i.kill()
		}
		osExit(3)
		return
		break
	default:
		fmt.Printf("Got a signal that wasn't handled. Please file a bug against bazel-watceer that describes how you did this. This is a big problem.\n")
	}

	i.interruptCount += 1
	if i.interruptCount > 2 {
		fmt.Printf("\nExiting from getting SIGINT 3 times\n")
		osExit(3)
	}
}

func (i *IBazel) newBazel() bazel.Bazel {
	b := bazelNew()
	b.SetArguments(i.bazelArgs)
	return b
}

func (i *IBazel) SetBazelArgs(args []string) {
	i.bazelArgs = args
}

func (i *IBazel) Cleanup() {
	i.buildFileWatcher.Close()
	i.sourceFileWatcher.Close()
}

func (i *IBazel) setup() error {
	var err error
	// Even though we are going to recreate this when the query happens, create
	// the pointer we will use to refer to the watchers right now.
	i.buildFileWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	i.sourceFileWatcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	i.sourceEventHandler = NewSourceEventHandler(i.sourceFileWatcher)

	return nil
}

// Run the specified target (singular) in the IBazel loop.
func (i *IBazel) Run(target string, args []string) error {
	i.args = args
	return i.loop("run", i.run, []string{target})
}

// Build the specified targets in the IBazel loop.
func (i *IBazel) Build(targets ...string) error {
	return i.loop("build", i.build, targets)
}

// Test the specified targets in the IBazel loop.
func (i *IBazel) Test(targets ...string) error {
	return i.loop("test", i.test, targets)
}

func (i *IBazel) loop(command string, commandToRun func(...string), targets []string) error {
	joinedTargets := strings.Join(targets, " ")

	i.state = QUERY
	for {
		i.iteration(command, commandToRun, targets, joinedTargets)
	}

	return nil
}

func (i *IBazel) iteration(command string, commandToRun func(...string), targets []string, joinedTargets string) {
	fmt.Printf("State: %s\n", i.state)
	switch i.state {
	case WAIT:
		select {
		case <-i.sourceEventHandler.SourceFileEvents:
			fmt.Printf("Detected source change. Rebuilding...\n")
			i.state = DEBOUNCE_RUN
		case <-i.buildFileWatcher.Events:
			fmt.Printf("Detected build graph change. Requerying...\n")
			i.state = DEBOUNCE_QUERY
		}
	case DEBOUNCE_QUERY:
		select {
		case <-i.buildFileWatcher.Events:
			i.state = DEBOUNCE_QUERY
		case <-time.After(debounceDuration):
			i.state = QUERY
		}
	case QUERY:
		// Query for which files to watch.
		fmt.Printf("Querying for BUILD files...\n")
		i.watchFiles(fmt.Sprintf(buildQuery, joinedTargets), i.buildFileWatcher)
		fmt.Printf("Querying for source files...\n")
		i.watchFiles(fmt.Sprintf(sourceQuery, joinedTargets), i.sourceFileWatcher)
		i.state = RUN
	case DEBOUNCE_RUN:
		select {
		case <-i.sourceEventHandler.SourceFileEvents:
			i.state = DEBOUNCE_RUN
		case <-time.After(debounceDuration):
			i.state = RUN
		}
	case RUN:
		i.state = WAIT
		fmt.Printf("%sing %s\n", strings.Title(command), joinedTargets)
		commandToRun(targets...)
	}
}

func (i *IBazel) build(targets ...string) {
	b := i.newBazel()

	b.Cancel()
	b.WriteToStderr(true)
	b.WriteToStdout(true)
	err := b.Build(targets...)
	if err != nil {
		fmt.Printf("Build error: %v", err)
		return
	}
}

func (i *IBazel) test(targets ...string) {
	b := i.newBazel()

	b.Cancel()
	b.WriteToStderr(true)
	b.WriteToStdout(true)
	err := b.Test(targets...)
	if err != nil {
		fmt.Printf("Build error: %v", err)
		return
	}
}

func (i *IBazel) subprocessRunning() bool {
	if i.cmd == nil {
		return false
	}
	if i.cmd.Process == nil {
		return false
	}
	if i.cmd.ProcessState != nil {
		if i.cmd.ProcessState.Exited() {
			return false
		}
	}

	return true
}

func (i *IBazel) kill() {
	if !i.subprocessRunning() {
		return
	}

	// Kill it with fire by sending SIGKILL to the process PID which should
	// propagate down to any subprocesses in the PGID (Process Group ID). To
	// send to the PGID, send the signal to the negative of the process PID.
	// Normally I would do this by calling i.cmd.Process.Signal, but that
	// only goes to the PID not the PGID.
	syscall.Kill(-i.cmd.Process.Pid, syscall.SIGKILL)
	i.cmd.Wait()
	i.cmd = nil
}

func (i *IBazel) run(targets ...string) {
	i.kill()

	b := i.newBazel()

	b.Cancel()
	b.WriteToStderr(true)
	b.WriteToStdout(true)

	tmpfile, err := ioutil.TempFile("", "bazel_script_path")
	if err != nil {
		fmt.Print(err)
	}
	// Close the file so bazel can write over it
	if err := tmpfile.Close(); err != nil {
		fmt.Print(err)
	}

	// Start by building the binary
	b.Run(append([]string{"--script_path=" + tmpfile.Name()}, targets...)...)

	targetPath := tmpfile.Name()

	// Now that we have built the target, construct a executable form of it for
	// execution in a go routine.
	i.cmd = execCommand(targetPath, i.args...)
	i.cmd.Stdout = os.Stdout
	i.cmd.Stderr = os.Stderr

	// Set a process group id (PGID) on the subprocess. This is
	i.cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Start run in a goroutine so that it doesn't block watching for files that
	// have changed.
	if err := i.cmd.Start(); err != nil {
		fmt.Printf("Error starting process: %v\n", err)
	}
}

func (i *IBazel) queryForSourceFiles(query string) []string {
	b := i.newBazel()
	b.WriteToStderr(false)
	b.WriteToStdout(false)

	res, err := b.Query(query)
	if err != nil {
		fmt.Printf("Error running Bazel %s\n", err)
	}

	toWatch := make([]string, 0, 10000)
	for _, line := range res {
		if strings.HasPrefix(line, "@") {
			continue
		}
		if strings.HasPrefix(line, "//external") {
			continue
		}

		// For files that are served from the root they will being with "//:". This
		// is a problematic string because, for example, "//:demo.sh" will become
		// "/demo.sh" which is in the root of the filesystem and is unlikely to exist.
		if strings.HasPrefix(line, "//:") {
			line = line[3:]
		}

		toWatch = append(toWatch, strings.Replace(strings.TrimPrefix(line, "//"), ":", "/", 1))
	}

	return toWatch
}

func (i *IBazel) watchFiles(query string, watcher *fsnotify.Watcher) {
	toWatch := i.queryForSourceFiles(query)

	// TODO: Figure out how to unwatch files that are no longer included

	for _, line := range toWatch {
		fmt.Printf("Watching: %s\n", line)
		err := watcher.Add(line)
		if err != nil {
			fmt.Printf("Error watching file %v\nError: %v\n", line, err)
			continue
		}
	}
}
