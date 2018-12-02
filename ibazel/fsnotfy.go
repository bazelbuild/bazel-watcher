package main

import (
	"github.com/fsnotify/fsnotify"
)

type fSNotifyWatcher interface {
	Close() error
	Add(name string) error
	Remove(name string) error
	Events() chan fsnotify.Event
	Errors() chan error
	Watcher() *fsnotify.Watcher
}

type realFSNotifyWatcher struct {
	w *fsnotify.Watcher
}

var _ fSNotifyWatcher = &realFSNotifyWatcher{}

func (w *realFSNotifyWatcher) Close() error                { return w.w.Close() }
func (w *realFSNotifyWatcher) Add(name string) error       { return w.w.Add(name) }
func (w *realFSNotifyWatcher) Remove(name string) error    { return w.w.Remove(name) }
func (w *realFSNotifyWatcher) Events() chan fsnotify.Event { return w.w.Events }
func (w *realFSNotifyWatcher) Errors() chan error          { return w.w.Errors }
func (w *realFSNotifyWatcher) Watcher() *fsnotify.Watcher  { return w.w }

func wrapWatcher(w *fsnotify.Watcher, err error) (fSNotifyWatcher, error) {
	return &realFSNotifyWatcher{w: w}, err
}
