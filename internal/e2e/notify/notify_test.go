package notify

import (
	"testing"

	"github.com/bazelbuild/bazel-watcher/internal/e2e"
	"github.com/bazelbuild/rules_go/go/tools/bazel_testing"
)

const mainFiles = `
-- BUILD --
sh_binary(
    name = "crash_script",
    srcs = ["crash_script.sh"],
)

sh_binary(
    name = "crash_script_notify",
    srcs = ["crash_script_notify.sh"],
    data = ["message.txt"],
    tags = [
        "ibazel_notify_changes",
    ],
)

-- crash_script.sh --
#!/bin/bash

echo "Starting crash script"
sleep 2
echo "Exiting after 2 seconds"
exit 1

-- message.txt --
HELLO

-- crash_script_notify.sh --
#!/bin/bash

# list all files in current directory
ls -l

# Filename to monitor
filename="message.txt"

# Read the initial content of the file
if [[ ! -f "$filename" ]]; then
  echo "Error: File $filename does not exist."
  exit 1
fi

content=$(cat "$filename")

# Function to check for CRASH in the file
check_for_crash() {
  if grep -q "CRASH" "$filename"; then
    echo "CRASH detected in $filename"
    exit 1
  fi
}

# Print the current content of the file
print_content() {
  echo "Message: $filename:"
  echo "$content"
}

# Monitor stdin for trigger message
while true; do
  print_content
  read -t 5 input
  if [[ $? -eq 0 ]]; then
    # Input received
    if [[ "$input" == "IBAZEL_BUILD_COMPLETED SUCCESS" ]]; then
      # Re-read the file
      if [[ ! -f "$filename" ]]; then
        echo "Error: File $filename does not exist."
        continue
      fi
      content=$(cat "$filename")
      check_for_crash
    fi
  fi
  sleep 1
done
`

func TestMain(m *testing.M) {
	bazel_testing.TestMain(m, bazel_testing.Args{
		Main: mainFiles,
	})
}

func TestNotifySubprocessCrashes(t *testing.T) {
	ibazel := e2e.SetUp(t)
	ibazel.Run([]string{}, "//:crash_script_notify")
	defer ibazel.Kill()

	// Expect the initial content of the file
	ibazel.ExpectOutput("Message: message.txt:")
	ibazel.ExpectOutput("HELLO")

	// It reacts to changes
	e2e.MustWriteFile(t, "message.txt", "WORLD")
	ibazel.ExpectOutput("WORLD")

	// This causes the subprocess to crash
	e2e.MustWriteFile(t, "message.txt", "CRASH")
	ibazel.ExpectOutput("CRASH detected in message.txt")

	// The subprocess is restarted, and the next change is detected
	e2e.MustWriteFile(t, "message.txt", "FIXED")
	ibazel.ExpectOutput("FIXED")
}
