package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"

	blaze_query "github.com/bazelbuild/bazel-watcher/third_party/bazel/master/src/main/protobuf"
	"github.com/golang/protobuf/proto"
)

func fuzzyEqual(pattern []interface{}, compare []string) bool {
	if len(pattern) != len(compare) {
		return false
	}

	for k, _ := range pattern {
		switch v := pattern[k].(type) {
		case *regexp.Regexp:
			if !v.MatchString(compare[k]) {
				return false
			}
			break
		case string:
			if v != compare[k] {
				return false
			}
			break
		default:
			fmt.Errorf("I have no idea what a %T is... please fix.", v)
			panic(v)
		}
	}

	return true
}

func writeLauncherScript(w io.Writer, launcher, pid string) (int, error) {
	script := `#! /usr/bin/env bash
set -e
%s &
PID="$!"
echo -n $PID > %s
wait`
	return fmt.Fprintf(w, script, launcher, pid)
}

func runWithScript(rawArgs []string) {
	// All run commands start with "run" as the first argument while all other
	// arguments are flags.
	// Translate ["run", "--script_path=/tmp/demo.sh", "//my/demo:target"] into
	// ["--script_path=/tmp/demo.sh", "//my/demo:target"]
	args := rawArgs[1:]

	fs := flag.NewFlagSet("bazel_run_flags", flag.ExitOnError)
	scriptPath := fs.String("script_path", "default", "")
	err := fs.Parse(args)
	if err != nil {
		panic(err)
	}

	sp, err := os.OpenFile(*scriptPath, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}

	err = sp.Chmod(0755)
	if err != nil {
		panic(err)
	}

	writeLauncherScript(sp, filepath.Join(os.TempDir(), "ibazel_e2e_subprocess_launcher"),
		filepath.Join(os.TempDir(), "ibazel_e2e_subprocess_launcher.pid"))

	err = sp.Close()
	if err != nil {
		panic(err)
	}
}

var sourceFiles = &blaze_query.QueryResult{
	Target: []*blaze_query.Target{
		&blaze_query.Target{
			Type: blaze_query.Target_SOURCE_FILE.Enum(),
			SourceFile: &blaze_query.SourceFile{
				Name:            proto.String("//e2e/simple:main.go"),
				Location:        proto.String("/home/user/go/src/github.com/bazelbuild/bazel-watcher/e2e/simple/BUILD.bazel:3:1"),
				VisibilityLabel: []string{"//visibility:private"},
			},
		},
	},
}
var buildFiles = &blaze_query.QueryResult{
	Target: []*blaze_query.Target{
		&blaze_query.Target{
			Type: blaze_query.Target_SOURCE_FILE.Enum(),
			SourceFile: &blaze_query.SourceFile{
				Name:            proto.String("//e2e/simple:BUILD.bazel"),
				Location:        proto.String("/home/user/go/src/github.com/bazelbuild/bazel-watcher/e2e/simple/BUILD.bazel:1"),
				VisibilityLabel: []string{"//visibility:private"},
			},
		},
	},
}

var target = &blaze_query.QueryResult{
	Target: []*blaze_query.Target{
		&blaze_query.Target{
			Type: blaze_query.Target_RULE.Enum(),
			Rule: &blaze_query.Rule{
				Name:      proto.String("//e2e/simple:simple"),
				RuleClass: proto.String("go_binary"),
				Attribute: []*blaze_query.Attribute{
					&blaze_query.Attribute{
						Name:            proto.String("name"),
						Type:            blaze_query.Attribute_STRING.Enum(),
						StringListValue: []string{"simple"},
					},
				},
				RuleInput: []string{"//e2e/simple:go_default_library"},
			},
		},
	},
}

func main() {
	inputs := []struct {
		args     []interface{}
		exitCode int
		output   interface{}
	}{
		{[]interface{}{"build"}, 0, `Build output`},

		// E2E simple test data.
		{[]interface{}{"query", "--output=proto", "--order_output=no", "buildfiles(deps(set(//e2e/simple)))"}, 0, buildFiles},
		{[]interface{}{"query", "--output=proto", "--order_output=no", "kind('source file', deps(set(//e2e/simple)))"}, 0, sourceFiles},
		{[]interface{}{"query", "--output=proto", "--order_output=no", "//e2e/simple"}, 0, target},
		{[]interface{}{"run", regexp.MustCompile("--script_path=.*"), "//e2e/simple"}, 0, runWithScript},
	}

	for _, opt := range inputs {
		if fuzzyEqual(opt.args, os.Args[1:]) {
			switch v := opt.output.(type) {
			case string:
				fmt.Println(opt.output)
				break
			case func([]string):
				v(os.Args[1:])
				break
			case *blaze_query.QueryResult:
				data, err := proto.Marshal(v)
				if err != nil {
					panic(err)
				}
				os.Stdout.Write(data)
			default:
				panic(fmt.Sprintf("Unkown output format %T", v))
			}
			os.Exit(opt.exitCode)
		}
	}

	logFile, err := os.OpenFile("/tmp/ibazel_test_run.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(logFile)
	}
	// Normally you would defer but I am calling os.Exit later which doesn't
	// call deferred functions.
	mw := io.MultiWriter(os.Stderr, logFile)
	fmt.Fprintf(mw, `Mock Bazel.

This is ALMOST CERTAINLY not useful to you. Please go to https://bazel.build
for more information on the Bazel project.

If you're interested in this tool, it is used for end to end testing of the
iBazel project.

Called with:
%v
`, os.Args)
	logFile.Sync()
	logFile.Close()
	os.Exit(255)
}
