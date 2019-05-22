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
