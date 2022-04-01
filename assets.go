//go:build !dev

package main

import (
	"embed"
	"io/fs"
	"path/filepath"
)

type embeddedAssets struct {
	embed.FS
}

func (e *embeddedAssets) Open(name string) (fs.File, error) {
	path := filepath.Join("web", "public", name)
	return e.FS.Open(path)
}

//go:embed web/public
var rootAssets embed.FS
var assets = &embeddedAssets{rootAssets}
