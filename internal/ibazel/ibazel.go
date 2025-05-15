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

package ibazel

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/bazelbuild/bazel-watcher/internal/bazel"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/command"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/fswatcher"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/fswatcher/common"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/lifecycle_hooks"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/live_reload"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/log"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/output_runner"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/profiler"
	"github.com/bazelbuild/bazel-watcher/internal/ibazel/workspace"
	"github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf/blaze_query"
)

var osExit = os.Exit
var bazelNew = bazel.New
var commandDefaultCommand = command.DefaultCommand
var commandNotifyCommand = command.NotifyCommand
var exitMessages = map[os.Signal]string{
	syscall.SIGINT:  "Subprocess killed from getting SIGINT (trigger SIGINT again to stop ibazel)",
	syscall.SIGTERM: "Subprocess killed from getting SIGTERM",
	syscall.SIGHUP:  "Subprocess killed from getting SIGHUP",
}

type State string
type runnableCommand func(...string) (*bytes.Buffer, error)

const (
	DEBOUNCE_QUERY State = "DEBOUNCE_QUERY"
	QUERY          State = "QUERY"
	WAIT           State = "WAIT"
	DEBOUNCE_RUN   State = "DEBOUNCE_RUN"
	RUN            State = "RUN"
	QUIT           State = "QUIT"
)

const sourceQuery = "kind('source file', deps(set(%s)))"
const buildQuery = "buildfiles(deps(set(%s)))"

type IBazel struct {
	debounceDuration time.Duration

	cmd         command.Command
	args        []string
	bazelArgs   []string
	startupArgs []string

	sigs           chan os.Signal // Signals channel for the current process
	interruptCount int

	workspaceFinder workspace.Workspace

	buildFileWatcher  common.Watcher
	sourceFileWatcher common.Watcher

	filesWatched map[common.Watcher]map[string]struct{} // Inner map is a surrogate for a set

	lifecycleListeners []Lifecycle

	state State
}

func New(version string) (*IBazel, error) {
	i := &IBazel{}
	err := i.setup()
	if err != nil {
		return nil, err
	}

	i.debounceDuration = 100 * time.Millisecond
	i.filesWatched = map[common.Watcher]map[string]struct{}{}
	i.workspaceFinder = &workspace.MainWorkspace{}

	i.sigs = make(chan os.Signal, 1)
	signal.Notify(i.sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	liveReload := live_reload.New()
	profiler := profiler.New(version)
	outputRunner := output_runner.New()
	lifecycleHooks := lifecycle_hooks.New()

	liveReload.AddEventsListener(profiler)

	i.lifecycleListeners = []Lifecycle{
		liveReload,
		profiler,
		outputRunner,
		lifecycleHooks,
	}

	info, _ := i.getInfo()
	for _, l := range i.lifecycleListeners {
		l.Initialize(&info)
	}

	go func() {
		for {
			i.handleSignals()
		}
	}()

	return i, nil
}

func (i *IBazel) handleSignals() {
	// Got an OS signal (SIGINT, SIGTERM, SIGHUP).
	sig := <-i.sigs

	if i.cmd == nil || !i.cmd.IsSubprocessRunning() {
		osExit(3)
		return
	}

	switch sig {
	case syscall.SIGINT:
		i.interruptCount++
		switch {
		case i.interruptCount > 2:
			log.NewLine()
			log.Fatal("Exiting from getting SIGINT 3 times")
			osExit(3)
		case i.interruptCount > 1:
			i.cmd.Kill()
		default:
			go func() {
				i.cmd.Terminate()
				log.NewLine()
				log.Log(exitMessages[sig])
			}()
		}
	case syscall.SIGTERM, syscall.SIGHUP:
		go func() {
			i.cmd.Terminate()
			log.NewLine()
			log.Log(exitMessages[sig])
			osExit(3)
		}()
	default:
		log.Fatal("Got a signal that wasn't handled. Please file a bug against bazel-watcher that describes how you did this. This is a big problem.")
	}
}

func (i *IBazel) newBazel() bazel.Bazel {
	b := bazelNew()
	b.SetStartupArgs(i.startupArgs)
	b.SetArguments(i.bazelArgs)
	return b
}

func (i *IBazel) SetBazelArgs(args []string) {
	i.bazelArgs = args
}

func (i *IBazel) SetStartupArgs(args []string) {
	i.startupArgs = args
}

func (i *IBazel) SetDebounceDuration(debounceDuration time.Duration) {
	i.debounceDuration = debounceDuration
}

func (i *IBazel) Cleanup() {
	i.buildFileWatcher.Close()
	i.sourceFileWatcher.Close()
	for _, l := range i.lifecycleListeners {
		l.Cleanup()
	}
}

func (i *IBazel) targetDecider(target string, rule *blaze_query.Rule) {
	for _, l := range i.lifecycleListeners {
		// TODO: As the name implies, it would be good to use this to make a
		// determination about if future events should be routed to this listener.
		// Why not do it now?
		// Right now I don't track which file is associated with the end target. I
		// just query for a list of all files that are rdeps of any target that is
		// in the list of targets to build/test/run (although run can only have 1).
		// Since I don't have that mapping right now the information doesn't
		// presently exist to implement this properly. Additionally, since querying
		// is currently in the critical path for getting something the user cares
		// about on screen, I'm not sure that it is wise to do this in the first
		// pass. It might be worth triggering the user action, launching their thing
		// and then running a background thread to access the data.
		l.TargetDecider(rule)
	}
}

func (i *IBazel) changeDetected(targets []string, changeType string, change string) {
	for _, l := range i.lifecycleListeners {
		l.ChangeDetected(targets, changeType, change)
	}
}

func (i *IBazel) beforeCommand(targets []string, command string) {
	for _, l := range i.lifecycleListeners {
		l.BeforeCommand(targets, command)
	}
}

func (i *IBazel) afterCommand(targets []string, command string, success bool, output *bytes.Buffer) {
	for _, l := range i.lifecycleListeners {
		l.AfterCommand(targets, command, success, output)
	}
}

func (i *IBazel) setup() error {
	var err error

	// Even though we are going to recreate this when the query happens, create
	// the pointer we will use to refer to the watchers right now.
	i.buildFileWatcher, err = fswatcher.NewWatcher()
	if err != nil {
		return err
	}

	i.sourceFileWatcher, err = fswatcher.NewWatcher()
	if err != nil {
		return err
	}

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

// Get code coverage for the specified targets in the IBazel loop.
func (i *IBazel) Coverage(targets ...string) error {
	return i.loop("coverage", i.test, targets)
}

func (i *IBazel) loop(command string, commandToRun runnableCommand, targets []string) error {
	joinedTargets := strings.Join(targets, " ")

	i.state = QUERY
	for {
		i.iteration(command, commandToRun, targets, joinedTargets)
	}
}

// fsnotify also triggers for file stat and read operations. Explicitly filter the modifying events
// to avoid triggering builds on file accesses (e.g. due to your IDE checking modified status).
const modifyingEvents = common.Write | common.Create | common.Rename | common.Remove

func (i *IBazel) iteration(command string, commandToRun runnableCommand, targets []string, joinedTargets string) {
	switch i.state {
	case WAIT:
		select {
		case e := <-i.sourceFileWatcher.Events():
			if _, ok := i.filesWatched[i.sourceFileWatcher][e.Name]; ok && e.Op&modifyingEvents != 0 {
				log.Logf("Changed: %q. Rebuilding...", e.Name)
				i.changeDetected(targets, "source", e.Name)
				i.state = DEBOUNCE_RUN
			}
		case e := <-i.buildFileWatcher.Events():
			if _, ok := i.filesWatched[i.buildFileWatcher][e.Name]; ok && e.Op&modifyingEvents != 0 {
				log.Logf("Build graph changed: %q. Requerying...", e.Name)
				i.changeDetected(targets, "graph", e.Name)
				i.state = DEBOUNCE_QUERY
			}
		}
	case DEBOUNCE_QUERY:
		select {
		case e := <-i.buildFileWatcher.Events():
			if _, ok := i.filesWatched[i.buildFileWatcher][e.Name]; ok && e.Op&modifyingEvents != 0 {
				i.changeDetected(targets, "graph", e.Name)
			}
			i.state = DEBOUNCE_QUERY
		case <-time.After(i.debounceDuration):
			i.state = QUERY
		}
	case QUERY:
		// Query for which files to watch.
		log.Logf("Querying for files to watch...")
		i.watchFiles(fmt.Sprintf(buildQuery, joinedTargets), i.buildFileWatcher)
		i.watchFiles(fmt.Sprintf(sourceQuery, joinedTargets), i.sourceFileWatcher)
		i.state = RUN
	case DEBOUNCE_RUN:
		select {
		case e := <-i.sourceFileWatcher.Events():
			if _, ok := i.filesWatched[i.sourceFileWatcher][e.Name]; ok && e.Op&modifyingEvents != 0 {
				i.changeDetected(targets, "source", e.Name)
			}
			i.state = DEBOUNCE_RUN
		case <-time.After(i.debounceDuration):
			i.state = RUN
		}
	case RUN:
		log.Logf("%s %s", strings.Title(verb(command)), joinedTargets)
		i.beforeCommand(targets, command)
		outputBuffer, err := commandToRun(targets...)
		i.interruptCount = 0
		i.afterCommand(targets, command, err == nil, outputBuffer)
		i.state = WAIT
	}
}

func verb(s string) string {
	switch s {
	case "run":
		return "running"
	default:
		return fmt.Sprintf("%sing", s)
	}
}

func (i *IBazel) build(targets ...string) (*bytes.Buffer, error) {
	b := i.newBazel()

	b.Cancel()
	b.WriteToStderr(true)
	b.WriteToStdout(true)
	outputBuffer, err := b.Build(targets...)
	if err != nil {
		log.Errorf("Build error: %v", err)
		return outputBuffer, err
	}
	return outputBuffer, nil
}

func (i *IBazel) test(targets ...string) (*bytes.Buffer, error) {
	b := i.newBazel()

	// Query the provided target patterns to construct a composite list
	// Make a set that represents all the found rules.
	targetRules := map[string]struct{}{}
	for _, target := range targets {
		rule, err := i.queryRule(target)
		if err != nil {
			log.Errorf("Error: %v", err)
		}
		targetRules[rule.GetName()] = struct{}{}
	}

	if len(targetRules) == 1 {
		setStream := true
		for _, arg := range b.Args() {
			if strings.HasPrefix(arg, "--test_output=") {
				setStream = false
			}
		}
		if setStream {
			log.Log("Found a single target test. Streaming results. You can override this by explicitly passing --test_output=summary")
			b.SetArguments(append([]string{"--test_output=streamed"}, b.Args()...))
		}
	}

	b.Cancel()
	b.WriteToStderr(true)
	b.WriteToStdout(true)
	outputBuffer, err := b.Test(targets...)
	if err != nil {
		log.Errorf("Build error: %v", err)
		return outputBuffer, err
	}
	return outputBuffer, err
}

func contains(l []string, e string) bool {
	for _, i := range l {
		if i == e {
			return true
		}
	}
	return false
}

func (i *IBazel) setupRun(target string) command.Command {
	rule, err := i.queryRule(target)
	if err != nil {
		log.Errorf("Error: %v", err)
	}

	i.targetDecider(target, rule)

	commandNotify := false
	for _, attr := range rule.Attribute {
		if *attr.Name == "tags" && *attr.Type == blaze_query.Attribute_STRING_LIST {
			if contains(attr.StringListValue, "ibazel_notify_changes") {
				commandNotify = true
			}
		}
	}

	if commandNotify {
		log.Logf("Launching with notifications")
		return commandNotifyCommand(i.startupArgs, i.bazelArgs, target, i.args)
	} else {
		return commandDefaultCommand(i.startupArgs, i.bazelArgs, target, i.args)
	}
}

func (i *IBazel) run(targets ...string) (*bytes.Buffer, error) {
	if i.cmd == nil {
		// If the command is empty, we are in our first pass through the state
		// machine and we need to make a command object.
		i.cmd = i.setupRun(targets[0])
		outputBuffer, err := i.cmd.Start()
		if err != nil {
			log.Errorf("Run start failed %v", err)
		}
		return outputBuffer, err
	}

	log.Logf("Notifying of changes")
	outputBuffer := i.cmd.NotifyOfChanges()
	return outputBuffer, nil
}

func (i *IBazel) queryRule(rule string) (*blaze_query.Rule, error) {
	b := i.newBazel()

	b.WriteToStderr(false)
	b.WriteToStdout(false)

	res, err := b.CQuery(i.cQueryArgs(rule)...)
	if err != nil {
		log.Errorf("Error running Bazel %v", err)
		osExit(4)
	}

	for _, target := range res.Results {
		switch *target.Target.Type {
		case blaze_query.Target_RULE:
			return target.Target.Rule, nil
		}
	}

	return nil, errors.New("No information available")
}

func (i *IBazel) getInfo() (map[string]string, error) {
	b := i.newBazel()

	res, err := b.Info()
	if err != nil {
		log.Errorf("Error getting Bazel info %v", err)
		return nil, err
	}

	return res, nil
}

func (i *IBazel) queryForSourceFiles(query string) ([]string, error) {
	b := i.newBazel()

	localRepositories, err := i.realLocalRepositoryPaths()
	if err != nil {
		return nil, err
	}

	res, err := b.Query(i.queryArgs(query)...)
	if err != nil {
		log.Errorf("Bazel query failed: %v", err)
		return nil, err
	}

	workspacePath, err := i.workspaceFinder.FindWorkspace()
	if err != nil {
		log.Errorf("Error finding workspace: %v", err)
		return nil, err
	}

	toWatch := make([]string, 0, len(res.GetTarget()))
	for _, target := range res.GetTarget() {
		switch *target.Type {
		case blaze_query.Target_SOURCE_FILE:
			label := target.GetSourceFile().GetName()
			if strings.HasPrefix(label, "@") {
				repo, target := parseTarget(label)
				if realPath, ok := findInLocalRepository(localRepositories, repo); ok {
					label = strings.Replace(target, ":", string(filepath.Separator), 1)
					toWatch = append(toWatch, filepath.Join(realPath, label))
					break
				}
				continue
			}
			if strings.HasPrefix(label, "//external") {
				continue
			}

			label = strings.Replace(strings.TrimPrefix(label, "//"), ":", string(filepath.Separator), 1)
			toWatch = append(toWatch, filepath.Join(workspacePath, label))
			break
		default:
			log.Errorf("%v\n", target)
		}
	}

	return toWatch, nil
}

func (i *IBazel) watchFiles(query string, watcher common.Watcher) {
	toWatch, err := i.queryForSourceFiles(query)
	if err != nil {
		// If the query fails, just keep watching the same files as before
		log.Errorf("Error querying for source files: %v", err)
		return
	}

	filesWatched := map[string]struct{}{}
	uniqueDirectories := map[string]struct{}{}

	for _, file := range toWatch {
		path, err := filepath.EvalSymlinks(file)
		if err != nil {
			log.Errorf("Error evaluating symbolic links for source file: %v", err)
			return
		}

		if _, err := os.Stat(path); !os.IsNotExist(err) {
			filesWatched[path] = struct{}{}
		}

		parentDirectory, _ := filepath.Split(path)
		// Add a watch to the file's parent directory, we might already have this dir in our set but thats OK
		uniqueDirectories[parentDirectory] = struct{}{}
	}

	watchList := keys(uniqueDirectories)
	err = watcher.UpdateAll(watchList)
	if err != nil {
		log.Errorf("Error(s) updating watch list:\n %v", err)
	}

	if len(filesWatched) == 0 {
		log.Errorf("Didn't find any files to watch from query %s", query)
	}

	i.filesWatched[watcher] = filesWatched
}

func (i *IBazel) queryArgs(args ...string) []string {
	queryArgs := append([]string(nil), args...)

	for _, arg := range i.bazelArgs {
		// List of args that should be passed to bazel query/cquery.
		if strings.HasPrefix(arg, "--override_repository=") {
			queryArgs = append(queryArgs, arg)
		}
	}

	return queryArgs
}

func (i *IBazel) cQueryArgs(args ...string) []string {
	// Unlike query, cquery can be affected by the majority of command line option.
	cQueryArgs := append([]string(nil), args...)
	cQueryArgs = append(cQueryArgs, i.bazelArgs...)
	return cQueryArgs
}

func parseTarget(label string) (repo string, target string) {
	parts := strings.Split(strings.TrimPrefix(label, "@"), "//")
	return parts[0], parts[1]
}

func keys(m map[string]struct{}) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func findInLocalRepository(localRepositories map[string]string, repo string) (string, bool) {
	if realPath, ok := localRepositories[repo]; ok {
		return realPath, ok
	}
	if realPath, ok := localRepositories[repo + "+"]; ok {
		return realPath, ok
	}
	return "", false
}

