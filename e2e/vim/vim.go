package simple

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/bazelbuild/bazel-watcher/e2e"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{})
}

func renameAndWriteNewFile(t *testing.T, fname, content string) {
	// write a file in the same manner as vim with backupcopy=no;
	// this will rename the original file to a file with a backup extension
	// and write the new file contents to the original filename

	fnameBackup := fmt.Sprintf("%s~", fname)
	e2e.Must(t, os.Rename(fname, fnameBackup))
	f, err := os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Write([]byte(content))
	if err != nil {
		t.Fatal(err)
	}
	e2e.Must(t, f.Close())
	e2e.Must(t, os.Remove(fnameBackup))
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
	e2e.Must(t, fBackup.Close())
	e2e.Must(t, f.Close())

	f, err = os.OpenFile(fname, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		t.Fatal(err)
	}
	_, err = f.Write([]byte(content))
	if err != nil {
		t.Fatal(err)
	}
	e2e.Must(t, f.Sync())
	e2e.Must(t, f.Close())
	e2e.Must(t, os.Remove(fnameBackup))
}

func TestSimpleRunWithModifiedFile_RenameAndWrite(t *testing.T) {
	e2e.MustMkdir(t, "vim")
	e2e.MustWriteFile(t, "vim/test.sh", `printf "Started!"`)
	e2e.MustWriteFile(t, "vim/BUILD", `
sh_binary(
	name = "test",
	srcs = ["test.sh"],
)
`)

	ibazel := e2e.NewIBazelTester(t)
	ibazel.Run([]string{}, "//vim:test")
	defer ibazel.Kill()
	ibazel.ExpectOutput("Started!")

	renameAndWriteNewFile(t, "vim/test.sh", `printf "Started2!"`)
	ibazel.ExpectOutput("Started2!")

	renameAndWriteNewFile(t, "vim/test.sh", `printf "Started3!"`)
	ibazel.ExpectOutput("Started3!")
}

func TestSimpleRunWithModifiedFile_CopyAndTruncWrite(t *testing.T) {
	e2e.MustMkdir(t, "truncate")
	e2e.MustWriteFile(t, "truncate/test.sh", `printf "Started!"`)
	e2e.MustWriteFile(t, "truncate/BUILD", `
sh_binary(
	name = "test",
	srcs = ["test.sh"],
)
`)

	ibazel := e2e.NewIBazelTester(t)
	ibazel.Run([]string{}, "//truncate:test")
	defer ibazel.Kill()
	ibazel.ExpectOutput("Started!")

	copyAndTruncWriteFile(t, "truncate/test.sh", `printf "Started2!"`)
	ibazel.ExpectOutput("Started2!")

	copyAndTruncWriteFile(t, "truncate/test.sh", `printf "Started3!"`)
	ibazel.ExpectOutput("Started3!")
}
