// Copyright 2017 The Bazel Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fsnotify


import (
	"errors"
	"runtime/debug"
	"testing"

	"github.com/fsnotify/fsnotify"
)

type mockFSNotifyWatcher struct {
	recentlyAddedFiles   map[string]struct{}
	recentlyRemovedFiles map[string]struct{}
	closed               bool
}
func (w *mockFSNotifyWatcher) Add(name string) error { 
	if _, ok := w.recentlyAddedFiles[name]; ok {
		return errors.New("Already added file " + name)
	}
	w.recentlyAddedFiles[name] = struct{}{}
	return nil
}
func (w *mockFSNotifyWatcher) Remove(name string) error { 
	if _, ok := w.recentlyRemovedFiles[name]; ok {
		return errors.New("Already removed file " + name)
	}
	w.recentlyRemovedFiles[name] = struct{}{}
	return nil
 }
func (w *mockFSNotifyWatcher) Close() error {
	if w.closed {
		return errors.New("Already closed")
	}
	w.closed = true
	return nil
}
func (w *mockFSNotifyWatcher) Events() chan fsnotify.Event {
	return nil
}

func (w *mockFSNotifyWatcher) Reset() {
	w.recentlyAddedFiles = make(map[string]struct{}, 0)
	w.recentlyRemovedFiles = make(map[string]struct{}, 0)
	w.closed = false
}

func (w *mockFSNotifyWatcher) assertRecentlyAdded(t *testing.T, added []string) {
	k := keys(w.recentlyAddedFiles)
	if val, ok := containsAll(k, added); !ok {
		t.Errorf("Expected Add(\"%s\") not to have been called", val)
		debug.PrintStack()
	}
	if val, ok := containsAll(added, k); !ok {
		t.Errorf("Expected Add(\"%s\") to have been called", val)
		debug.PrintStack()
	}
}

func (w *mockFSNotifyWatcher) assertRecentlyRemoved(t *testing.T, removed []string) {
	k := keys(w.recentlyRemovedFiles)
	if val, ok := containsAll(k, removed); !ok {
		t.Errorf("Expected Remove(\"%s\") not to have been called", val)
		debug.PrintStack()
	}
	if val, ok := containsAll(removed, k); !ok {
		t.Errorf("Expected Remove(\"%s\") to have been called", val)
		debug.PrintStack()
	}
}

func (w *mockFSNotifyWatcher) assertClosed(t *testing.T, closed bool) {
	if w.closed != closed {
		assertion := "to"
		if !closed {
			assertion = "not to"
		}
		t.Errorf("Expected Close() %s have been called", assertion)
		debug.PrintStack()
	}
}

func newWatcher() (*realFSNotifyWatcher, *mockFSNotifyWatcher) {
	mock := &mockFSNotifyWatcher{}
	mock.Reset()
	watcher := &realFSNotifyWatcher{wrapper: mock}
	return watcher, mock
}

func TestWatchedFilesState(t *testing.T) {
	watcher, mock := newWatcher()

	mock.Reset()
	watcher.UpdateAll([]string{
		"/path/a",
		"/path/b",
		"/path/c",
	})
	mock.assertRecentlyAdded(t, []string{
		"/path/a",
		"/path/b",
		"/path/c",
	})
	mock.assertRecentlyRemoved(t, []string{})
	
	mock.Reset()
	watcher.UpdateAll([]string{
		"/path/a",
		"/path/b",
	})
	mock.assertRecentlyAdded(t, []string{})
	mock.assertRecentlyRemoved(t, []string{
		"/path/c",
	})

	mock.Reset()
	watcher.UpdateAll([]string{
		"/path/a",
		"/path/d",
	})
	mock.assertRecentlyAdded(t, []string{
		"/path/d",
	})
	mock.assertRecentlyRemoved(t, []string{
		"/path/b",
	})

	mock.Reset()
	watcher.UpdateAll([]string{})
	mock.assertRecentlyAdded(t, []string{})
	mock.assertRecentlyRemoved(t, []string{
		"/path/a",
		"/path/d",
	})

	mock.Reset()
	watcher.UpdateAll([]string{
		"/other/1",
		"/other/2",
		"/other/4",
	})
	mock.assertRecentlyAdded(t, []string{
		"/other/1",
		"/other/2",
		"/other/4",
	})
	mock.assertRecentlyRemoved(t, []string{})


	mock.assertClosed(t, false)
	watcher.Close()
	mock.assertClosed(t, true)
}

// Equal tells whether a and b contain the same elements, regardless of order
func containsAll(a, b []string) (string, bool) {
	OUTER:
    for _, v1 := range a {
        for _, v2 := range b {
			if v1 == v2 {
				continue OUTER
			}
		}
		return v1, false
	}
    return "", true
}

func keys(m map[string]struct{}) []string {
	keys := make([]string, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
