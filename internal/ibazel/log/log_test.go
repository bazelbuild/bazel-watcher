package log

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const prefix = "github.com/bazelbuild/bazel-watcher/internal/ibazel/log."

func TestNonfLoggers(t *testing.T) {
	tests := []struct {
		method interface{}
		msg    string
		args   []interface{}
		want   string
		exits  bool
		color  color
	}{
		{
			method: Log,
			msg:    "log",
			want:   "log",
			exits:  false,
			color:  logColor,
		},
		{
			method: Logf,
			msg:    "log %d",
			args:   []interface{}{123},
			want:   "log 123",
			exits:  false,
			color:  logColor,
		},
		{
			method: Error,
			msg:    "error",
			want:   "error",
			exits:  false,
			color:  errorColor,
		},
		{
			method: Errorf,
			msg:    "error %d",
			args:   []interface{}{123},
			want:   "error 123",
			exits:  false,
			color:  errorColor,
		},
		{
			method: Fatal,
			msg:    "fatal",
			want:   "fatal",
			exits:  true,
			color:  fatalColor,
		},
		{
			method: Fatalf,
			msg:    "fatal %d",
			args:   []interface{}{123},
			want:   "fatal 123",
			exits:  true,
			color:  fatalColor,
		},
	}

	for _, test := range tests {
		funcName := strings.TrimPrefix(
			runtime.FuncForPC(
				reflect.ValueOf(test.method).Pointer()).Name(), prefix)

		t.Run(fmt.Sprintf("%v(%q, %v)", funcName, test.msg, test.args), func(t *testing.T) {
			calledExit := false
			osExit = func(int) {
				calledExit = true
			}
			timeNow = func() time.Time {
				parsedTime, err := time.Parse(time.RFC3339, "2019-11-13T00:05:07+00:00")
				if err != nil {
					t.Errorf("Couldn't parse time: %v", err)
				}
				return parsedTime
			}

			buf := &bytes.Buffer{}
			SetWriter(buf)

			switch f := test.method.(type) {
			case func(string):
				f(test.msg)
			case func(string, ...interface{}):
				f(test.msg, test.args...)
			}

			if test.exits && calledExit == false {
				t.Errorf("Fatal should call exit")
			}

			got := buf.String()
			want := fmt.Sprintf("%siBazel [12:05AM]\x1b[0m: %s\n", test.color, test.want)
			if diff := cmp.Diff(got, want); diff != "" {
				t.Errorf("\nGot:  %q\nWant: %q\nDiff:\n%s", got, want, diff)
			}

			buf.Reset()
		})
	}
}

func TestBanner(t *testing.T) {
	buf := &bytes.Buffer{}
	SetWriter(buf)

	Banner("This is multi", "line output that", "is expected to be printed")

	got := buf.String()
	want := fmt.Sprintf(`
%s################################################################################%s
%s#%s This is multi                                                                %s#%s
%s#%s line output that                                                             %s#%s
%s#%s is expected to be printed                                                    %s#%s
%s################################################################################%s

`, bannerColor, resetColor, bannerColor, resetColor, bannerColor, resetColor, bannerColor, resetColor, bannerColor, resetColor, bannerColor, resetColor, bannerColor, resetColor, bannerColor, resetColor)
	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("\nGot:  %q\nWant: %q\nDiff:\n%s", got, want, diff)
	}
}
