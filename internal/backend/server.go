package backend

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"kalimah/internal/backend/middleware"
	"math"
	"math/rand"
	"net/http"
	"sort"
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
	router.GET("/api/words/surah/:surah/page/:page", s.GetWords)
	router.GET("/api/tafsir/surah/:surah/ayah/:ayah", s.GetTafsir)
	router.POST("/api/track", s.TrackWord)

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

func (s *Server) GetWords(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var err error
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}()

	// Parse URL params
	page, _ := strconv.Atoi(ps.ByName("page"))
	surah, _ := strconv.Atoi(ps.ByName("surah"))

	// Prepare read only transaction
	tx, err := s.DB.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()

	// Fetch ayah count for this surah
	var nAyah int
	err = tx.Get(&nAyah, `SELECT end-start+1 FROM surah WHERE id = ?`, surah)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	// Adjust pagination
	nAyahPerPage := 30
	maxPage := int(math.Ceil(float64(nAyah) / float64(nAyahPerPage)))
	if page < 1 {
		page = 1
	} else if page > maxPage {
		page = maxPage
	}

	// Fetch words for this page
	words := []Word{}
	err = tx.Select(&words,
		`WITH last_word AS (
			SELECT IFNULL(last_word, 0) id FROM tracker WHERE id = 1),
		ayah_range AS (
			SELECT (30*(?-1) + 1) start, MIN(end, 30*?) end 
			FROM surah
			WHERE id = ?)
		SELECT w.id, w.ayah-ar.start+1 ayah, w.position,
			w.arabic, w.translation, w.id <= lw.id answered,
			w.ayah <> LEAD(w.ayah, 1, w.ayah+1) OVER (ORDER BY w.ayah) is_separator
		FROM word w, ayah_range ar, last_word lw
		WHERE w.ayah >= ar.start AND w.ayah <= ar.end
		ORDER BY w.id`, page, page, surah)
	if err != nil && err != sql.ErrNoRows {
		return
	}

	// Fetch choice candidate
	var choiceCandidates []string
	err = tx.Select(&choiceCandidates,
		`SELECT DISTINCT w.translation FROM word w
		ORDER BY RANDOM() LIMIT ?`, len(words)*5)
	if err != nil {
		return
	}

	// Apply choice to each word
	nCandidates := len(choiceCandidates)
	for i, word := range words {
		// Prepare choices for this word
		choices := make([]Choice, 10)
		choices[0] = Choice{Text: word.Translation, IsCorrect: true}

		// Fetch incorrect choice randomly
		usedCandidateIdx := map[int]struct{}{}
		for j := 0; j < 9; j++ {
			var candidateIdx int

			// Make sure candidate is unused and not correct
			for {
				candidateIdx = rand.Intn(nCandidates)
				_, candidateIsUsed := usedCandidateIdx[candidateIdx]
				candidateIsCorrect := choiceCandidates[candidateIdx] == word.Translation
				if !candidateIsUsed && !candidateIsCorrect {
					break
				}
			}

			// Save the candidate
			usedCandidateIdx[candidateIdx] = struct{}{}
			choices[j+1] = Choice{Text: choiceCandidates[candidateIdx], IsCorrect: false}
		}

		// Sort the choices
		sort.Slice(choices, func(i, j int) bool {
			return choices[i].Text < choices[j].Text
		})

		// Apply choices to word
		words[i].Choices = choices
	}

	// Create return data
	data := struct {
		CurrentPage int    `json:"currentPage"`
		MaxPage     int    `json:"maxPage"`
		Words       []Word `json:"words"`
	}{
		CurrentPage: page,
		MaxPage:     maxPage,
		Words:       words,
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(&data)
}

func (s *Server) GetTafsir(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

func (s *Server) TrackWord(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Prepare error handling
	var err error
	defer func() {
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	}()

	// Decode response
	var currentWord Word
	err = json.NewDecoder(r.Body).Decode(&currentWord)
	if err != nil {
		return
	}

	// Update tracker
	_, err = s.DB.Exec(
		"UPDATE tracker SET last_word = ? WHERE id = 1",
		currentWord.ID)
	if err != nil {
		return
	}
}
