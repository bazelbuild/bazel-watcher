package main

import (
	"github.com/fsnotify/fsnotify"
)

type Watcher interface {
	Cleanup()
	WatchSourceFiles(files []string) (count int, err error)
	WatchBuildFiles(files []string) (count int, err error)

	BuildEvents() chan fsnotify.Event
	SourceEvents() chan fsnotify.Event
}
