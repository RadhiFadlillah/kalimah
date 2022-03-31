package backend

import (
	"fmt"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

var presetMimeTypes = map[string]string{
	".css":  "text/css; charset=utf-8",
	".html": "text/html; charset=utf-8",
	".js":   "application/javascript",
	".png":  "image/png",
}

// ServeAssets serve the assets for specified file path.
func ServeAssets(w http.ResponseWriter, r *http.Request, assets fs.FS, filePath string) error {
	// Get request header
	reqEtag := r.Header.Get("If-None-Match")
	reqLastModified := r.Header.Get("If-Modified-Since")

	// Open file
	src, err := assets.Open(filePath)
	if err != nil {
		return err
	}
	defer src.Close()

	// Get file statistic
	info, err := src.Stat()
	if err != nil {
		return err
	}

	// Check if file is modified
	etag := fmt.Sprintf(`W/"%x-%x"`, info.ModTime().Unix(), info.Size())
	lastModified := info.ModTime().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	if reqEtag == etag || reqLastModified == lastModified {
		w.WriteHeader(http.StatusNotModified)
		return nil
	}

	// Get content type
	ext := filepath.Ext(filePath)
	mimeType := guessTypeByExtension(ext)

	// Write response header
	w.Header().Set("ETag", etag)
	w.Header().Set("Last-Modified", lastModified)
	w.Header().Set("Content-Length", strconv.FormatInt(info.Size(), 10))

	if mimeType != "" {
		w.Header().Set("Content-Type", mimeType)
		w.Header().Set("X-Content-Type-Options", "nosniff")
	}

	// Serve file
	_, err = io.Copy(w, src)
	return err
}

func guessTypeByExtension(ext string) string {
	ext = strings.ToLower(ext)
	if v, ok := presetMimeTypes[ext]; ok {
		return v
	}

	return mime.TypeByExtension(ext)
}
