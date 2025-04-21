package ui

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed assets/js/* assets/css/*
var assetFS embed.FS

func Assets(strip string) http.Handler {
	sub, err := fs.Sub(assetFS, "assets")
	if err != nil {
		panic(err)
	}
	return http.StripPrefix(strip, http.FileServerFS(sub))
}
