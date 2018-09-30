package simple

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"testing"

	bazel "github.com/bazelbuild/bazel-integration-testing/go"
	"github.com/bazelbuild/bazel-watcher/e2e"
)

func must(t *testing.T, e error) {
	if e != nil {
		t.Errorf("Error: %s", e)
		debug.PrintStack()
	}
}

func TestSimpleRun(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchFileWithMode("test.sh", `printf "Started!"`, 0777))
	must(t, b.ScratchFile("BUILD", `
sh_binary(
	name = "test",
	srcs = ["test.sh"],
)
`))

	ibazel := e2e.NewIBazelTester(t, b)
	ibazel.Run("//:test")
	defer ibazel.Kill()

	ibazel.ExpectOutput("Started!")
}

func TestSimpleRunUnderSubdir(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchDir("subdir"))
	must(t, b.ScratchFileWithMode("subdir/test.sh", `printf "Started!"`, 0777))
	must(t, b.ScratchFile("subdir/BUILD", `
sh_binary(
	name = "test",
	srcs = ["test.sh"],
)
`))

	ibazel := e2e.NewIBazelTester(t, b)

	err = os.Chdir("subdir")
	if err != nil {
		t.Fatal(err)
	}

	ibazel.Run("test")
	defer ibazel.Kill()

	err = os.Chdir("..")
	if err != nil {
		t.Fatal(err)
	}

	ibazel.ExpectOutput("Started!")
}

func TestSimpleRunWithModifiedFile(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchFileWithMode("test.sh", `printf "Started!"`, 0777))
	must(t, b.ScratchFile("BUILD", `
sh_binary(
	name = "test",
	srcs = ["test.sh"],
)
`))

	ibazel := e2e.NewIBazelTester(t, b)
	ibazel.Run("//:test")
	defer ibazel.Kill()

	// Give it time to start up and query.
	must(t, b.ScratchFileWithMode("test.sh", `printf "Started2!"`, 0777))
	ibazel.ExpectOutput("Started2!")

	// Manipulate a source file and sleep past the debounce.
	must(t, b.ScratchFileWithMode("test.sh", `printf "Started3!"`, 0777))
	ibazel.ExpectOutput("Started3!")

	// Now a BUILD file.
	must(t, b.ScratchFile("BUILD", `
sh_binary(
	# New comment
	name = "test",
	srcs = ["test.sh"],
)
`))
	ibazel.ExpectOutput("Started3!")
}

func renameAndWriteNewFile(t *testing.T, fname, content string) {
	// write a file in the same manner as vim with backupcopy=no;
	// this will rename the original file to a file with a backup extension
	// and write the new file contents to the original filename

	fnameBackup := fmt.Sprintf("%s~", fname)
	must(t, os.Rename(fname, fnameBackup))
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Write([]byte(content))
	if err != nil {
		t.Fatal(err)
	}
	must(t, f.Close())
	must(t, os.Remove(fnameBackup))
}

func copyAndTruncWriteFile(t *testing.T, fname string, content string) {
	// write a file in the same manner as vim with backupcopy=yes;
	// this will copy the file to a suffixed backup file,
	// truncate the existing file and write the new content
	fnameBackup := fmt.Sprintf("%s~", fname)

	f, err := os.Open(fname)
	if err != nil {
		t.Fatal(err)
	}
	fBackup, err := os.OpenFile(fnameBackup, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(fBackup, f)
	if err != nil {
		t.Fatal(err)
	}
	must(t, fBackup.Close())
	must(t, f.Close())

	f, err = os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Write([]byte(content))
	if err != nil {
		t.Fatal(err)
	}
	must(t, f.Sync())
	must(t, f.Close())
	must(t, os.Remove(fnameBackup))
}

func TestSimpleRunWithModifiedFile_RenameAndWrite(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchFileWithMode("test.sh", `printf "Started!"`, 0777))
	must(t, b.ScratchFile("BUILD", `
sh_binary(
	name = "test",
	srcs = ["test.sh"],
)
`))

	ibazel := e2e.NewIBazelTester(t, b)
	ibazel.Run("//:test")
	defer ibazel.Kill()
	ibazel.ExpectOutput("Started!")

	renameAndWriteNewFile(t, "test.sh", `printf "Started2!"`)
	ibazel.ExpectOutput("Started2!")

	renameAndWriteNewFile(t, "test.sh", `printf "Started3!"`)
	ibazel.ExpectOutput("Started3!")
}

func TestSimpleRunWithModifiedFile_CopyAndTruncWrite(t *testing.T) {
	b, err := bazel.New()
	if err != nil {
		t.Fatal(err)
	}
	must(t, b.ScratchFile("WORKSPACE", ""))
	must(t, b.ScratchFileWithMode("test.sh", `printf "Started!"`, 0777))
	must(t, b.ScratchFile("BUILD", `
sh_binary(
	name = "test",
	srcs = ["test.sh"],
)
`))

	ibazel := e2e.NewIBazelTester(t, b)
	ibazel.Run("//:test")
	defer ibazel.Kill()
	ibazel.ExpectOutput("Started!")

	copyAndTruncWriteFile(t, "test.sh", `printf "Started2!"`)
	ibazel.ExpectOutput("Started2!")

	copyAndTruncWriteFile(t, "test.sh", `printf "Started3!"`)
	ibazel.ExpectOutput("Started3!")
}
