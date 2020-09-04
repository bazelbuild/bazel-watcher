---
name: Bug report
about: Create a report for ibazel
title: ''
labels: ''
assignees: ''

---

**Describe the bug**
A clear and concise description of what the bug is.

**Reproduction instructions**
A link to a repo where this can be tested

Steps to reproduce the behavior:
1. `git clone https://github.com/pathto/repo`
2. `ibazel test //path/to:target`
3. Edit `//path/to/file.go` to have a smiley face at the top
4. See error

**Expected behavior**
A clear and concise description of what you expected to happen.

**`bazel query --output=build //path/to:target`**
If your issue is about a specific target in your system working, Please include the output of `bazel query --output=build //path/to:target`,

**Version (please complete the following information):**
 - OS: [e.g. iOS]
 - Browser [e.g. chrome, safari]
 - ibazel Version [e.g. 0.2.4] (run `ibazel 2>&1 | head -n 1` to get this)
 - Bazel version [e.g. 3.1.0] (run `bazel version` to get this. Please include all the lines)
 
**Additional context**
Add any other context about the problem here.
