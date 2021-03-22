#! /usr/bin/env python3

# Note that this file should work in both Python 2 and 3.

from __future__ import print_function
from subprocess import Popen, PIPE

dirty = Popen(["git", "diff-index", "--quiet", "HEAD"], stdout=PIPE).wait() != 0

commit_process = Popen(["git", "describe", "--always", "--tags", "--abbrev=0"], stdout=PIPE)
(version, err) = commit_process.communicate()

print("STABLE_GIT_VERSION %s%s" % (
    version.decode("utf-8").replace("\n", ""),
    "-dirty" if dirty else "")
)

