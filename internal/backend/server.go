package backend

import (
	"database/sql"
	"encoding/json"
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
	router.GET("/api/surah", s.GetSurah)
	router.GET("/api/word", s.GetWords)
	router.GET("/api/answer", s.GetAnswers)
	router.POST("/api/answer", s.SubmitAnswer)

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
	err := serveAssets(w, r, s.Assets, r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func (s *Server) GetSurah(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error
	defer markHttpError(w, err)

	listSurah := []Surah{}
	err = s.DB.Select(&listSurah,
		`WITH last_word AS (
			SELECT IFNULL(MAX(last_word), 0) id FROM tracker WHERE id = 1),
		translated_surah AS (
			SELECT DISTINCT surah id
			FROM word, last_word
			WHERE word.id <= last_word.id)
		SELECT s.id, s.name, s.translation, IIF(ts.id IS NULL, 0, 1) translated
		FROM surah s
		LEFT JOIN translated_surah ts ON s.id = ts.id
		ORDER BY s.id`)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&listSurah)
}

func (s *Server) GetWords(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error
	defer markHttpError(w, err)

	strSurah := r.URL.Query().Get("surah")
	surah, _ := strconv.Atoi(strSurah)

	words := []Word{}
	err = s.DB.Select(&words,
		`WITH last_word AS (
			SELECT IFNULL(MAX(last_word), 0) id FROM tracker WHERE id = 1)
		SELECT w.id, ayah, position, arabic,
			IIF(w.id <= lw.id, translation, '') translation
		FROM word w, last_word lw
		WHERE surah = ?
		ORDER BY w.id, ayah, position`, surah)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&words)
}

func (s *Server) GetAnswers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error
	defer markHttpError(w, err)

	strWord := r.URL.Query().Get("word")
	word, _ := strconv.Atoi(strWord)

	answers := []string{}
	err = s.DB.Select(&answers,
		`WITH correct_answer AS (
			SELECT translation FROM word w WHERE id = ?),
		wrong_answer AS (
			SELECT DISTINCT w.translation 
			FROM word w, correct_answer cw
			WHERE w.translation NOT LIKE "%"||cw.translation||"%"
			ORDER BY RANDOM()
			LIMIT 9),
		all_answer AS (
			SELECT * FROM correct_answer
			UNION ALL
			SELECT * FROM wrong_answer)
		SELECT translation answer FROM all_answer 
		ORDER BY random()`, word)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&answers)
}

func (s *Server) SubmitAnswer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Prepare error handling
	var err error
	defer markHttpError(w, err)

	// Decode response
	var answer Answer
	err = json.NewDecoder(r.Body).Decode(&answer)
	if err != nil {
		return
	}

	// Compare with database
	var correctAnswer string
	err = s.DB.Get(&correctAnswer,
		"SELECT translation FROM word WHERE id = ?",
		answer.WordID)
	if err != nil {
		return
	}

	// Update tracker
	var responseCode int
	if answer.Answer == correctAnswer {
		responseCode = 1
		_, err = s.DB.Exec(
			"UPDATE tracker SET last_word = ? WHERE id = 1",
			answer.WordID)
		if err != nil {
			return
		}
	}

	w.Header().Add("Content-Type", "text/plain")
	_, err = fmt.Fprint(w, responseCode)
}
