//go:build dev

package main

import (
	"io/fs"
	"net/http"
)

type embeddedAssets struct {
	http.Dir
}

func (e *embeddedAssets) Open(name string) (fs.File, error) {
	return e.Dir.Open(name)
}

var assets = &embeddedAssets{http.Dir("web/public")}
