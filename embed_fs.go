//go:build !test

package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed ui/dist/*
var embeddedFiles embed.FS

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(embeddedFiles, "ui/dist")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
