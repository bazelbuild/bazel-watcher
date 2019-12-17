package log

import (
	"fmt"
	"io"
	"os"
	"time"
)

type color string

var writer io.Writer = os.Stderr
var osExit = os.Exit
var timeNow = time.Now

const (
	reset color = "\033[0m"

	errorColor color = "\033[31m"
	fatalColor color = "\033[41m"
	logColor   color = "\033[96m"
)

func log(c color, msg string, args ...interface{}) {
	fmt.Fprintf(writer, "%siBazel [%s]%s: ",
		c,
		timeNow().Local().Format(time.Kitchen),
		reset)
	fmt.Fprintf(writer, msg, args...)
	fmt.Fprintf(writer, "\n")
}

// NewLine prints a new line to the screen without any preamble.
func NewLine() {
	fmt.Fprintf(writer, "\n")
}

// Error prints an error to the screen with a preamble.
func Error(msg string) {
	Errorf(msg)
}

// Errorf prints an error to the screen with a preamble.
func Errorf(msg string, args ...interface{}) {
	log(errorColor, msg, args...)
}

// Fatal prints a fatal error to the screen with a preamble.
func Fatal(msg string) {
	Fatalf(msg)
}

// Fatalf prints a fatal error to the screen with a preamble.
func Fatalf(msg string, args ...interface{}) {
	log(fatalColor, msg, args...)
	osExit(1)
}

// Log prints a message to the screen with a preamble.
func Log(msg string) {
	Logf(msg)
}

// Logf prints a message to the screen with a preamble.
func Logf(msg string, args ...interface{}) {
	log(logColor, msg, args...)
}

// SetWriter decides which io.Writer to write logs to.
func SetWriter(w io.Writer) {
	writer = w
}

// FakeExit makes the Fatal log methods not exit.
func FakeExit() {
	osExit = func(int) {}
}
