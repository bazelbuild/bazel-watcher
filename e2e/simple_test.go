package e2e

import (
	"testing"
	"time"
)

func TestSimpleRun(t *testing.T) {
	t.Skip()
	ibazel := IBazelTester("//e2e/simple", "e2e/simple/simple")
	ibazel.Run()
	defer ibazel.Kill()
	time.Sleep(10 * time.Millisecond)
	res := ibazel.GetOutput()

	assertEqual(t, "Started!", res, "Ouput was inequal")
}

func TestSimpleRunWithModifiedFile(t *testing.T) {
	ibazel := IBazelTester("//e2e/simple", "e2e/simple/simple")
	ibazel.Run()
	defer ibazel.Kill()

	pids := []int64{}
	count := 0
	verify := func() {
		time.Sleep(500 * time.Millisecond)

		count += 1

		pid := ibazel.GetSubprocessPid()
		for _, v := range pids {
			if pid == v {
				t.Errorf("Subsequent runs of the subcommand should have differing pids. %v, %v", pid, v)
			}
		}
		pids = append(pids, pid)

		expectedOut := ""
		for i := 0; i < count; i++ {
			expectedOut += "Started!"
		}
		assertEqual(t, expectedOut, ibazel.GetOutput(), "Ouput was inequal")
	}

	// Give it time to start up and query.
	verify()

	// Manipulate a source file and sleep past the debounce.
	manipulateSourceFile(count)
	verify()

	// Now a BUILD file.
	manipulateBUILDFile(count)
	verify()
}
