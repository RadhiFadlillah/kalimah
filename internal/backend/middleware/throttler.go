package middleware

import (
	"net/http"
	"strings"
	"time"
)

// Throttler is middleware to throttle response speed. Used to emulate low connection speed.
// Throttler will only triggered if it occured in URL prefixed with `/api/`
type Throttler struct {
	Handler http.Handler
	Delay   time.Duration
}

func NewThrottler(handler http.Handler, delay time.Duration) http.Handler {
	return &Throttler{handler, delay}
}

func (t *Throttler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// If filter defined and this URL pass the filter, delay the response
	if strings.HasPrefix(r.URL.Path, "/api/") {
		if t.Delay == 0 {
			t.Delay = 500 * time.Millisecond
		}
		time.Sleep(t.Delay)
	}

	t.Handler.ServeHTTP(w, r)
}
