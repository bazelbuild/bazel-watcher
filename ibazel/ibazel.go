package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bazelbuild/bazel-watcher/bazel"
	"github.com/fsnotify/fsnotify"
)

var bazelNew = bazel.New

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

	return i, nil
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
func (i *IBazel) Run(target string) error {
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
	b := bazelNew()

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
	b := bazelNew()

	b.Cancel()
	b.WriteToStderr(true)
	b.WriteToStdout(true)
	err := b.Test(targets...)
	if err != nil {
		fmt.Printf("Build error: %v", err)
		return
	}
}

func (i *IBazel) run(targets ...string) {
	b := bazelNew()

	b.Cancel()
	b.WriteToStderr(true)
	b.WriteToStdout(true)

	// Start run in a goroutine so that it doesn't block watching for files that
	// have changed.
	go b.Run(targets...)
}

func queryForSourceFiles(query string) []string {
	b := bazelNew()
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
	toWatch := queryForSourceFiles(query)

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
