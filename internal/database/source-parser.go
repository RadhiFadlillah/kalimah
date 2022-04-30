package database

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"embed"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/yuin/goldmark"
)

type SurahRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type SurahTranslation struct {
	Name        string `json:"name"`
	Translation string `json:"translation"`
}

type Word struct {
	Ayah     int    `json:"ayah"`
	Position int    `json:"position"`
	Uthmani  string `json:"uthmani"`
	Nastaliq string `json:"nastaliq"`
}

var (
	//go:embed source
	sourceAssets embed.FS
	rxTafsirAyah = regexp.MustCompile(`^=+\s*(\d+)\s*=+$`)
)

func parseSurahRange() (map[int]SurahRange, error) {
	// Open source
	f, err := sourceAssets.Open("source/surah.json.gz")
	if err != nil {
		return nil, fmt.Errorf("open failed: %w", err)
	}
	defer f.Close()

	// Decompress data
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("extract failed: %w", err)
	}
	defer gz.Close()

	// Decode data
	data := map[int]SurahRange{}
	err = json.NewDecoder(gz).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("decode JSON failed: %w", err)
	}

	return data, nil
}

func parseSurahTranslation() (map[int]SurahTranslation, error) {
	// Open source
	f, err := sourceAssets.Open("source/surah-indonesia.json.gz")
	if err != nil {
		return nil, fmt.Errorf("open failed: %w", err)
	}
	defer f.Close()

	// Decompress data
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("extract failed: %w", err)
	}
	defer gz.Close()

	// Decode data
	data := map[int]SurahTranslation{}
	err = json.NewDecoder(gz).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("decode JSON failed: %w", err)
	}

	return data, nil
}

func parseAyahTranslation() (map[int]string, error) {
	// Open source
	f, err := sourceAssets.Open("source/ayah-indonesia.json.gz")
	if err != nil {
		return nil, fmt.Errorf("open failed: %w", err)
	}
	defer f.Close()

	// Decompress data
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("extract failed: %w", err)
	}
	defer gz.Close()

	// Decode translation data
	var data struct {
		Translations map[int]string `json:"translations"`
		Footnotes    map[int]string `json:"footnotes"`
	}

	err = json.NewDecoder(gz).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("decode JSON failed: %w", err)
	}

	return data.Translations, nil
}

func parseAyahTafsir() (map[int]string, error) {
	// Open source
	f, err := sourceAssets.Open("source/ayah-tafsir-indonesia.md.gz")
	if err != nil {
		return nil, fmt.Errorf("open failed: %w", err)
	}
	defer f.Close()

	// Decompress data
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("extract failed: %w", err)
	}
	defer gz.Close()

	// Scan each line
	var currentID int
	var currentContent []string
	tafsirs := map[int]string{}
	scanner := bufio.NewScanner(gz)

	fnSave := func() error {
		if len(currentContent) == 0 || currentID == 0 {
			return nil
		}

		var buf bytes.Buffer
		content := strings.Join(currentContent, "\n\n")
		err = goldmark.Convert([]byte(content), &buf)
		if err != nil {
			return err
		}

		content = strings.TrimSpace(buf.String())
		tafsirs[currentID] = content
		return nil
	}

	for scanner.Scan() {
		// Fetch the current line
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if it's ayah number
		numberParts := rxTafsirAyah.FindStringSubmatch(line)
		if len(numberParts) > 0 {
			if err = fnSave(); err != nil {
				return nil, err
			}

			currentID, _ = strconv.Atoi(numberParts[1])
			currentContent = []string{}
			continue
		}

		// Save the line
		currentContent = append(currentContent, line)
	}

	// Save the trailing content
	if err = fnSave(); err != nil {
		return nil, err
	}

	return tafsirs, nil
}

func parseWord() (map[int]Word, error) {
	// Open source
	f, err := sourceAssets.Open("source/word.json.gz")
	if err != nil {
		return nil, fmt.Errorf("open failed: %w", err)
	}
	defer f.Close()

	// Decompress data
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("extract failed: %w", err)
	}
	defer gz.Close()

	// Decode data
	data := map[int]Word{}
	err = json.NewDecoder(gz).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("decode JSON failed: %w", err)
	}

	return data, nil
}

func parseWordTranslation() (map[int]string, error) {
	// Open source
	f, err := sourceAssets.Open("source/word-indonesia.json.gz")
	if err != nil {
		return nil, fmt.Errorf("open failed: %w", err)
	}
	defer f.Close()

	// Decompress data
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("extract failed: %w", err)
	}
	defer gz.Close()

	// Decode data
	data := map[int]string{}
	err = json.NewDecoder(gz).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("decode JSON failed: %w", err)
	}

	return data, nil
}
