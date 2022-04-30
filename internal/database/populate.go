package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	unicodeTransformer = transform.Chain(norm.NFKD, norm.NFKC)
)

func PopulateData(db *sqlx.DB) error {
	// Create transaction
	logrus.Println("opening transaction")
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	// If error ever happened, rollback
	defer func() {
		if err != nil && tx != nil {
			tx.Rollback()
		}
	}()

	// Populate data
	logrus.Println("populate surah")
	if err = populateSurah(tx); err != nil {
		return fmt.Errorf("failed to populate surah: %v", err)
	}

	logrus.Println("populate ayah")
	if err = populateAyah(tx); err != nil {
		return fmt.Errorf("failed to populate ayah: %v", err)
	}

	logrus.Println("populate words")
	if err = populateWord(tx); err != nil {
		return fmt.Errorf("failed to populate word: %v", err)
	}

	logrus.Println("populate tracker")
	_, err = tx.Exec(`
		INSERT INTO tracker (id) VALUES (1) 
		ON CONFLICT DO NOTHING`)
	if err != nil {
		return fmt.Errorf("failed to populate tracker: %v", err)
	}

	// Commit transaction
	logrus.Println("commit transaction")
	err = tx.Commit()

	return nil
}

func populateSurah(tx *sqlx.Tx) error {
	// Parse surah
	ranges, err := parseSurahRange()
	if err != nil {
		return err
	}

	translations, err := parseSurahTranslation()
	if err != nil {
		return err
	}

	// Prepare query statement
	stmt, err := tx.Preparex(`
		INSERT INTO surah (id, name, translation, start, end)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT DO UPDATE
		SET name = excluded.name,
			translation = excluded.translation,
			start = excluded.start,
			end = excluded.end`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute queries
	for id := 1; id <= 114; id++ {
		trans := translations[id]
		surahRange := ranges[id]

		_, err = stmt.Exec(id, trans.Name, trans.Translation, surahRange.Start, surahRange.End)
		if err != nil {
			return err
		}
	}

	return nil
}

func populateAyah(tx *sqlx.Tx) error {
	// Open data
	translations, err := parseAyahTranslation()
	if err != nil {
		return err
	}

	tafsirs, err := parseAyahTafsir()
	if err != nil {
		return err
	}

	// Prepare query statement
	stmt, err := tx.Preparex(`
		INSERT INTO ayah (id, translation, tafsir)
		VALUES (?, ?, ?)
		ON CONFLICT DO UPDATE
		SET translation = excluded.translation,
			tafsir = excluded.tafsir`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute queries
	for id := 1; id <= 6236; id++ {
		tafsir := tafsirs[id]
		translation := translations[id]
		_, err = stmt.Exec(id, translation, tafsir)
		if err != nil {
			return err
		}
	}

	return nil
}

func populateWord(tx *sqlx.Tx) error {
	// Open data
	words, err := parseWord()
	if err != nil {
		return err
	}

	translations, err := parseWordTranslation()
	if err != nil {
		return err
	}

	// Prepare query statement
	stmt, err := tx.Preparex(`
		INSERT INTO word (id, ayah, position, arabic, translation)
		VALUES (?, ?, ?, ?, ?)
		ON CONFLICT DO UPDATE
		SET ayah = excluded.ayah,
			position = excluded.position,
			arabic = excluded.arabic,
			translation = excluded.translation`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute queries
	for id := 1; id <= len(words); id++ {
		word := words[id]
		translation := translations[id]
		arabic := normalizeUnicode(word.Nastaliq)

		_, err = stmt.Exec(id, word.Ayah, word.Position, arabic, translation)
		if err != nil {
			return err
		}
	}

	return nil
}

func normalizeUnicode(str string) string {
	normalized, _, err := transform.String(unicodeTransformer, str)
	if err != nil {
		return str
	}
	return normalized
}
