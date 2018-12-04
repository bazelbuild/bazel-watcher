# Bazel watcher

[![Build status](https://badge.buildkite.com/7694a2e22dcb7ea2e2ec80bb7e8e0380c700079e761394096f.svg?branch=master)](https://buildkite.com/bazel/bazel-watcher-postsubmit)

Note: This is not an official Google product.

A source file watcher for [Bazel](https://Bazel.build) projects

Ever wanted to save a file and have your tests automatically run? How about
restart your webserver when one of the source files change? Look no further.

Install `ibazel` using one of the 3 methods [described below](#installation). Then:

```bash
# ibazel build //path/to/my:target
```

Hack hack hack. Save and your target will be rebuilt.

Right now this repo supports `build`, `test`, and `run`.

## Installation

There are currently 3 ways to install iBazel

### Mac (Homebrew)

If you run a mac you can install it from [homebrew](https://brew.sh).

```
$ brew tap bazelbuild/tap
$ brew tap-pin bazelbuild/tap
$ brew install ibazel
```

### NPM

If you're a JS developer you can install it as a `devDependency` or by calling `npm install` directly in your project

```
npm install @bazel/ibazel
```

### Compiling yourself

You can, of course, build iBazel using Bazel.

```
git clone git@github.com:bazelbuild/bazel-watcher
cd bazel-watcher
bazel build //ibazel
```

Now copy the generated binary onto your path:

```bash
export PATH=$PATH:$PWD/bazel-bin/ibazel/$GOOS_$GOARCH_pure_stripped
```

where `$GOOS` and `$GOARCH` are your host OS (e.g., `darwin` or `linux`) and architecture (e.g., `amd64`).

## Running a target

By default, a target started with `ibazel run` will be terminated and restarted
whenever it's notified of source changes. Alternatively, if the build rule for
your target contains `ibazel_notify_changes` in its `tags` attribute, then the
command will stay alive and will receive a notification of the source changes on
stdin.

## Profiling

iBazel has a `--profile_dev` flag which turns on a generated profile output file
about the build process and timing. To use it include this flag in the command line. For example,

```
iBazel --profile_dev=profile.json run devserver
```

The profile outfile is in concatenated JSON format. Events are outputted one per line.

### Profiler events

| Event | Description | Attributes <font size=1>(* optional)</font> |
| ------------- | ------------- | ------------- |
| `IBAZEL_START` | Emitted when iBazel is started as part of the first iteration | `type`, `iteration`, `time`, `iBazelVersion`, `bazelVersion`, `maxHeapSize`, `committedHeapSize` |
| `SOURCE_CHANGE` | A source file change was detected | `type`, `iteration`, `time`, `targets`, `elapsed`, `change` |
| `GRAPH_CHANGE` | A build file change was detected | `type`, `iteration`, `time`, `targets`, `elapsed`, `changes`* |
| `RELOAD_TRIGGERED` | A livereload was triggered to any listening browsers | `type`, `iteration`, `time`, `targets`, `elapsed`, `changes`* |
| `RUN_START` | A run operation started | `type`, `iteration`, `time`, `targets`, `elapsed`, `changes`* |
| `RUN_FAILED` | A run operation failed | `type`, `iteration`, `time`, `targets`, `elapsed`, `changes`* |
| `RUN_DONE` | A run operation completed successfully | `type`, `iteration`, `time`, `targets`, `elapsed`, `changes`* |
| `BUILD_START` | A build operation started | `type`, `iteration`, `time`, `targets`, `elapsed`, `changes`* |
| `BUILD_FAILED` | A build operation failed | `type`, `iteration`, `time`, `targets`, `elapsed`, `changes`* |
| `BUILD_DONE` | A build operation completed successfully | `type`, `iteration`, `time`, `targets`, `elapsed`, `changes`* |
| `TEST_START` | A test operation started | `type`, `iteration`, `time`, `targets`, `elapsed`, `changes`* |
| `TEST_FAILED` | A test operation failed | `type`, `iteration`, `time`, `targets`, `elapsed`, `changes`* |
| `TEST_DONE` | A test operation completed successfully | `type`, `iteration`, `time`, `targets`, `elapsed`, `changes`* |
| `REMOTE_EVENT` | A remote event was received from the browser | `type`, `iteration`, `time`, `targets`, `elapsed`, `remoteType`, `remoteTime`, `remoteElapsed`, `remoteData` |
| `REMOTE_EVENT / PAGE_LOAD` | A remote event emitted by the profiler client-side script on the browser's `load` event. `remoteType` is `PAGE_LOAD`. | `type`, `iteration`, `time`, `targets`, `elapsed`, `remoteType`, `remoteTime`, `remoteElapsed`, `remoteData` |

### Event attributes

| Attribute | Type | Description |
| ------------- | ------------- | ------------- |
| `type` | string | Event type. |
| `iteration` | string | Unique build iteration key. |
| `time` | integer | Time of event. |
| `targets` | string[] | List of targets that are being built (Note: this is a complete list and includes targets that were already built prior to an iteration). |
| `elapsed` | integer | Elapsed time in ms since the start of the iteration. |
| `change` | string | The file changed on a `SOURCE_CHANGE` or `GRAPH_CHANGE` event. |
| `changes` | string[] | A cumulative list of files changed during a build iteration. |
| `iBazelVersion` | string | Version of iBazel that generated this event. |
| `bazelVersion` | string | Version of bazel in use. |
| `maxHeapSize` | string | Max heap size as reported by Bazel. |
| `committedHeapSize` | string | Committed heap size as reported by Bazel. |
| `remoteType` | string | Sub-type for `REMOTE_EVENT` type. |
| `remoteTime` | number | Browser time for `REMOTE_EVENT` type. |
| `remoteElapsed` | number | Elapsed time in browser since `navigationStart` for `REMOTE_EVENT` type. |
| `remoteData` | string | Data sent from browser for `REMOTE_EVENT` type. This may be in escaped JSON format for some remote events. |

### Example profile output file

You can find an example profile output JSON file [here](https://github.com/bazelbuild/bazel-watcher/blob/master/example.profile.json). Below is the file in pretty print JSON format:

```
{  
   "type":"IBAZEL_START",
   "iteration":"4214114686684e0f",
   "time":1513706108351,
   "iBazelVersion":"v0.2.0-dirty",
   "bazelVersion":"release 0.8.1-homebrew",
   "maxHeapSize":"3817MB",
   "committedHeapSize":"1372MB"
}
{  
   "type":"RUN_START",
   "iteration":"4214114686684e0f",
   "time":1513706109329,
   "targets":["//src:devserver"],
   "elapsed":978
}
{  
   "type":"RELOAD_TRIGGERED",
   "iteration":"4214114686684e0f",
   "time":1513706114595,
   "targets":["//src:devserver"],
   "elapsed":6244
}
{  
   "type":"RUN_DONE",
   "iteration":"4214114686684e0f",
   "time":1513706114595,
   "targets":["//src:devserver"],
   "elapsed":6244
}
{  
   "type":"SOURCE_CHANGE",
   "iteration":"7e6f8e150e9a8367",
   "time":1513706129384,
   "targets":["//src:devserver"],
   "change":"/Users/greg/google/gregmagolan/angular-bazel-example/src/hello-world/hello-world.component.ts"
}
{  
   "type":"RUN_START",
   "iteration":"7e6f8e150e9a8367",
   "time":1513706129484,
   "targets":["//src:devserver"],
   "elapsed":100,
   "changes":["/Users/greg/google/gregmagolan/angular-bazel-example/src/hello-world/hello-world.component.ts"]
}
{  
   "type":"RELOAD_TRIGGERED",
   "iteration":"7e6f8e150e9a8367",
   "time":1513706133947,
   "targets":["//src:devserver"],
   "elapsed":4563,
   "changes":["/Users/greg/google/gregmagolan/angular-bazel-example/src/hello-world/hello-world.component.ts"]
}
{  
   "type":"RUN_DONE",
   "iteration":"7e6f8e150e9a8367",
   "time":1513706133947,
   "targets":["//src:devserver"],
   "elapsed":4563,
   "changes":["/Users/greg/google/gregmagolan/angular-bazel-example/src/hello-world/hello-world.component.ts"]
}
{  
   "type":"REMOTE_EVENT",
   "iteration":"7e6f8e150e9a8367",
   "time":1513706134297,
   "targets":["//src:devserver"],
   "elapsed":4913,
   "remoteType":"PAGE_LOAD",
   "remoteTime":1513706134294,
   "remoteElapsed":346,
   "remoteData":"{\"pageLoadTime\":344,\"fetchTime\":9,\"connectTime\":0,\"requestTime\":6,\"responseTime\":6,\"renderTime\":325,\"navigationStart\":1513706133948,\"unloadEventStart\":1513706133962,\"unloadEventEnd\":1513706133962,\"redirectStart\":0,\"redirectEnd\":0,\"fetchStart\":1513706133952,\"domainLookupStart\":1513706133952,\"domainLookupEnd\":1513706133952,\"connectStart\":1513706133952,\"connectEnd\":1513706133952,\"secureConnectionStart\":0,\"requestStart\":1513706133955,\"responseStart\":1513706133955,\"responseEnd\":1513706133961,\"domLoading\":1513706133967,\"domInteractive\":1513706134222,\"domContentLoadedEventStart\":1513706134222,\"domContentLoadedEventEnd\":1513706134222,\"domComplete\":1513706134292,\"loadEventStart\":1513706134292}"
}
```

## Remote events

Remote events require the client-side profiling script. If you are using the `ts_devserver` bazel rule, this script will automatically be included in the development bundle so you don't have to worry about including it. If you're not using `ts_devserver` for development mode, you can include the following script tag to pull in the client-side profiling script:

```
<script src="http://localhost:30000/profiler.js"></script>
```

The script is served by iBazel on port 30000 by default. If port 30000 is not available, iBazel will attempt to find another available port between 30001 and 30099.

Remote events in the profiler script are sent using the [Beacon API](https://developer.mozilla.org/en-US/docs/Web/API/Beacon_API). This API is available in Chrome 39, Firefox 31, Edge and Opera 26. It is not available in Internet Explorer or Safari. Browser compatability can be checked [here](https://developer.mozilla.org/en-US/docs/Web/API/Navigator/sendBeacon#Browser_compatibility).

If your browser does not support the Beacon API, you will see the following error in the console when including the profiler script: `iBazel profiler disabled because Beacon API is not available`.

## Custom remote events

When the profiler script is loaded, a `window.IBazelProfileEvent(eventType, data)` public API is made available for generating custom remote events. This function sends a custom REMOTE_EVENT to the iBazel profiler log.

| Param | Type | Description |
| ------------- | ------------- | ------------- |
| `eventType` | string | The event type that ends up in the 'remoteType' attribute of the REMOTE_EVENT. |
| `data` | any | Optional data associated with the event. This is converted to a string. If it is an object it will be converted to escaped JSON in the profiler log. |

## Additional notes

### Termination

SIGINT has to be sent twice to kill ibazel: once to kill the subprocess, and
the second time for ibazel itself. Also, ibazel will exit on its own when a
bazel query fails, but it will stay alive when a build, test, or run fails.
We use an exit code of 3 for a signal termination, and 4 for a query failure.
These codes are not an API and may change at any point.

### What about the `--watchfs` flag?

Bazel has a flag called `--watchfs` which, according to the bazel command-line
help does:

> If true, Bazel tries to use the operating system's file watch service for
> local changes instead of scanning every file for a change

Unfortunately, this option does not rebuild the project on save like the Bazel
watcher does, but instead queries the file system for a list of files that have
been invalidated since last build and will require reinspection by the Bazel
server.

### Big thanks

 * [Google](http://opensource.google.com) for cross-platform build/test CI instances.
 * [Sauce Labs](https://saucelabs.com) for cross-browser testing.

Copyright 2017 The Bazel Authors. All right reserved.
