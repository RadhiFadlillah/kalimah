package database

import (
	"fmt"
	"net/url"
	"time"

	"github.com/jmoiron/sqlx"
)

// Open database on specified path.
func Open(dbPath string) (db *sqlx.DB, err error) {
	// Prepare DSN
	q := url.Values{}
	q.Add("_foreign_keys", "1")
	q.Add("_journal_mode", "WAL")
	dsn := dbPath + "?" + q.Encode()

	// Connect to database
	db, err = sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}
	db.SetConnMaxLifetime(time.Minute)

	// Create transaction
	var tx *sqlx.Tx
	tx, err = db.Beginx()
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	// If error ever happened, rollback and close database
	defer func() {
		if err != nil {
			if tx != nil {
				tx.Rollback()
			}
			if db != nil {
				db.Close()
			}
			db = nil
		}
	}()

	// Generate tables
	ddlQueries := []string{
		ddlCreateSurah,
		ddlCreateAyah,
		ddlCreateWord,
		ddlCreateTracker}

	for _, query := range ddlQueries {
		_, err = tx.Exec(query)
		if err != nil {
			return nil, err
		}
	}

	// Commit transaction
	err = tx.Commit()
	return db, err
}
