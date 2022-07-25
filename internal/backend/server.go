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
	router.GET("/api/surah/:surah/word", s.GetSurahWords)
	router.GET("/api/surah/:surah/ayah/:ayah", s.GetSurahAyah)
	router.GET("/api/choice/:word-id", s.GetChoices)
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
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}()

	listSurah := []Surah{}
	err = s.DB.Select(&listSurah,
		`WITH last_word AS (
			SELECT IFNULL(last_word, 0) + 1 id FROM tracker WHERE id = 1),
		last_ayah AS (
			SELECT MAX(word.ayah) ayah
			FROM word, last_word
			WHERE word.id <= last_word.id)
		SELECT s.id, s.name, s.translation, (s.start <= la.ayah) translated
		FROM surah s, last_ayah la
		ORDER BY s.id`)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&listSurah)
}

func (s *Server) GetSurahWords(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}()

	surah, _ := strconv.Atoi(ps.ByName("surah"))

	words := []Word{}
	err = s.DB.Select(&words,
		`WITH last_word AS (
			SELECT IFNULL(last_word, 0) id FROM tracker WHERE id = 1),
		ayah_range AS (
			SELECT start, end FROM surah WHERE id = ?)
		SELECT w.id, w.ayah-ar.start+1 ayah, w.position, w.arabic,
			IIF(w.id <= lw.id, translation, '') translation,
			w.ayah <> LEAD(w.ayah, 1, w.ayah+1) OVER (ORDER BY w.ayah) is_separator
		FROM word w, ayah_range ar, last_word lw
		WHERE w.ayah >= ar.start AND w.ayah <= ar.end
		ORDER BY w.id`, surah)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&words)
}

func (s *Server) GetSurahAyah(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}()

	// Parse parameter
	surah, _ := strconv.Atoi(ps.ByName("surah"))
	ayah, _ := strconv.Atoi(ps.ByName("ayah"))

	// Fetch translation and tafsir
	var data Ayah
	err = s.DB.Get(&data,
		`SELECT id, translation, tafsir FROM ayah
		WHERE id = ? - 1 + (SELECT start FROM surah WHERE id = ?)`,
		ayah, surah)
	if err != nil {
		return
	}

	// Fetch arabic text
	err = s.DB.Get(&data.Arabic,
		`SELECT GROUP_CONCAT(arabic, " ") FROM word 
		WHERE ayah = ? - 1 + (SELECT start FROM surah WHERE id = ?)
		ORDER BY position`,
		ayah, surah)
	if err != nil {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&data)
}

func (s *Server) GetChoices(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}()

	strWordId := ps.ByName("word-id")
	wordId, _ := strconv.Atoi(strWordId)

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
		ORDER BY LOWER(answer)`, wordId)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&answers)
}

func (s *Server) SubmitAnswer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Prepare error handling
	var err error
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}()

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
		answer.Word.ID)
	if err != nil {
		return
	}

	// Update tracker
	var responseCode int
	if answer.Answer == correctAnswer {
		responseCode = 1

		if answer.Word.IsSeparator {
			_, err = s.DB.Exec(
				"UPDATE tracker SET last_word = ? WHERE id = 1",
				answer.Word.ID)
			if err != nil {
				return
			}
		}
	}

	w.Header().Add("Content-Type", "text/plain")
	_, err = fmt.Fprint(w, responseCode)
}
