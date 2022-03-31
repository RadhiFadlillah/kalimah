package middleware

import (
	"net/http"

	"github.com/NYTimes/gziphandler"
)

// NewGzipper returns a middleware to GZip all of HTTP response.
func NewGzipper(handler http.Handler) http.Handler {
	return gziphandler.GzipHandler(handler)
}
