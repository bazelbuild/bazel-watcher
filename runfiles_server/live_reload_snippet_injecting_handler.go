// Copyright 2019 The Bazel Authors. All rights reserved.
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

package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
)

type liveReloadSnippetInjectingHandler struct {
	http.Handler
	snippet []byte
}

func (l *liveReloadSnippetInjectingHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if !l.shouldInjectLiveReloadSnippet(req) {
		l.Handler.ServeHTTP(resp, req)
		return
	}
	responseWriter := &responseWriter{
		ResponseWriter:                    resp,
		liveReloadSnippetInjectingHandler: l,
	}
	l.Handler.ServeHTTP(responseWriter, req)
	if _, err := resp.Write(l.snippet); err != nil {
		log.Printf("warning: %v", err)
	}
}

func (l *liveReloadSnippetInjectingHandler) shouldInjectLiveReloadSnippet(req *http.Request) bool {
	return l.snippet != nil && strings.HasSuffix(req.URL.Path, ".html")
}

type responseWriter struct {
	http.ResponseWriter
	*liveReloadSnippetInjectingHandler
}

func (d *responseWriter) WriteHeader(statusCode int) {
	contentLength := d.ResponseWriter.Header().Get("Content-Length")
	var length, err = strconv.Atoi(contentLength)
	if err != nil {
		log.Printf("couldn't parse Content-Length header %v", contentLength)
	} else {
		length += len(d.liveReloadSnippetInjectingHandler.snippet)
		d.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(length))
	}
	d.ResponseWriter.WriteHeader(statusCode)
}
