package log

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type color string

var osExit = os.Exit
var timeNow = time.Now

type writerLogger struct {
	w io.Writer
}

func NewWriterLogger(w io.Writer) Logger {
	return &writerLogger{w}
}

func (w *writerLogger) Log(a ...interface{}) {
	fmt.Fprint(w.w, a...)
	fmt.Fprint(w.w, "\n")
}

var logger Logger = &writerLogger{os.Stderr}

type Logger interface{ Log(...interface{}) }

const (
	resetColor  color = "\033[0m"
	bannerColor color = "\033[33m"
	errorColor  color = "\033[31m"
	fatalColor  color = "\033[41m"
	logColor    color = "\033[96m"
)

func log(c color, msg string, args ...interface{}) {
	if t, ok := logger.(interface {
		Helper()
	}); ok {
		t.Helper()
	}

	userMsg := fmt.Sprintf(msg, args...)
	logger.Log(fmt.Sprintf("%siBazel [%s]%s: %s",
		c,
		timeNow().Local().Format(time.Kitchen),
		resetColor,
		userMsg))
}

// NewLine prints a new line to the screen without any preamble.
func NewLine() {
	if t, ok := logger.(interface {
		Helper()
	}); ok {
		t.Helper()
	}

	logger.Log("")
}

// Print out a banner surrounded by # to draw attention to the eye.
func Banner(lines ...string) {
	if t, ok := logger.(interface {
		Helper()
	}); ok {
		t.Helper()
	}

	NewLine()
	logger.Log(fmt.Sprintf("%s%s%s", bannerColor, strings.Repeat("#", 80), resetColor))

	for _, line := range lines {
		logger.Log(fmt.Sprintf("%s#%s %-76s %s#%s", bannerColor, resetColor, line, bannerColor, resetColor))
	}

	logger.Log(fmt.Sprintf("%s%s%s", bannerColor, strings.Repeat("#", 80), resetColor))
	NewLine()
}

// Error prints an error to the screen with a preamble.
func Error(msg string) {
	if t, ok := logger.(interface {
		Helper()
	}); ok {
		t.Helper()
	}

	Errorf(msg)
}

// Errorf prints an error to the screen with a preamble.
func Errorf(msg string, args ...interface{}) {
	if t, ok := logger.(interface {
		Helper()
	}); ok {
		t.Helper()
	}

	log(errorColor, msg, args...)
}

// Fatal prints a fatal error to the screen with a preamble.
func Fatal(msg string) {
	if t, ok := logger.(interface {
		Helper()
	}); ok {
		t.Helper()
	}

	Fatalf(msg)
}

// Fatalf prints a fatal error to the screen with a preamble.
func Fatalf(msg string, args ...interface{}) {
	if t, ok := logger.(interface {
		Helper()
	}); ok {
		t.Helper()
	}

	log(fatalColor, msg, args...)
	exit(1)
}

// Log prints a message to the screen with a preamble.
func Log(msg string) {
	if t, ok := logger.(interface {
		Helper()
	}); ok {
		t.Helper()
	}
	Logf(msg)
}

// Logf prints a message to the screen with a preamble.
func Logf(msg string, args ...interface{}) {
	if t, ok := logger.(interface {
		Helper()
	}); ok {
		t.Helper()
	}

	log(logColor, msg, args...)
}

func SetTesting(t interface {
	Log(...interface{})
	Fail()
}) {
	logger = t
}

// SetLogger decides which io.Writer to write logs to.
func SetLogger(l Logger) {
	logger = l
}

// FakeExit makes the Fatal log methods not exit.
func FakeExit() {
	osExit = func(int) {}
}

func exit(code int) {
	osExit(code)
}
