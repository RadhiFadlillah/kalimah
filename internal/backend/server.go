package backend

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"kalimah/internal/backend/middleware"
	"net/http"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

// Server is server for serving app.
type Server struct {
	DB      *sqlx.DB
	Assets  fs.FS
	DevMode bool
}

// Serve serves app in specified port.
func (s *Server) Serve(port int) error {
	// Create router
	router := httprouter.New()
	router.GET("/", s.ServeIndex)
	router.GET("/res/*filepath", s.ServeFile)
	router.GET("/build/*filepath", s.ServeFile)

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, arg interface{}) {
		http.Error(w, fmt.Sprintf("unrecovered error: %v", arg), 500)
	}

	// Apply middlewares
	var handler http.Handler = router
	handler = middleware.NewGzipper(handler)

	if s.DevMode {
		handler = middleware.NewThrottler(handler, 500*time.Millisecond)
	}

	// Create server
	url := fmt.Sprintf(":%d", port)
	svr := &http.Server{
		Addr:         url,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: time.Minute,
	}

	// Serve app
	logrus.Println("serve app in", url)
	return svr.ListenAndServe()
}

// ServeIndex serves the index page for admin app.
func (s *Server) ServeIndex(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Open HTML file
	f, err := s.Assets.Open("app.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer f.Close()

	// Read the entire content
	bt, err := ioutil.ReadAll(f)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// Serve file
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(int(len(bt))))

	_, err = w.Write(bt)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// ServeFile serves the files.
func (s *Server) ServeFile(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	err := ServeAssets(w, r, s.Assets, r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
