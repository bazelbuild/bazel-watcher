module github.com/bazelbuild/bazel-watcher

require (
	github.com/bazelbuild/rules_go v0.53.0
	github.com/fsnotify/fsevents v0.1.1
	github.com/fsnotify/fsnotify v1.6.0
	github.com/golang/protobuf v1.5.4
	github.com/google/go-cmp v0.5.9
	github.com/gorilla/websocket v1.5.0
	github.com/jaschaephraim/lrserver v0.0.0-20171129202958-50d19f603f71
	github.com/mattn/go-shellwords v1.0.12
	github.com/smartystreets/goconvey v1.7.2 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
)

require golang.org/x/sys v0.28.0

require google.golang.org/protobuf v1.36.3 // indirect

go 1.22

toolchain go1.24.1
