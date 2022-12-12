//go:build test

package main

import (
	"net/http"
	"testing/fstest"
)

func getFileSystem() http.FileSystem {
	fsys := fstest.MapFS{
		"index.html": {},
	}
	return http.FS(fsys)
}
