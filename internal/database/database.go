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
	q.Add("_pragma", "foreign_keys=1")
	q.Add("_pragma", "synchronous=0")
	dsn := dbPath + "?" + q.Encode()

	// Connect to database
	db, err = sqlx.Connect("sqlite", dsn)
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
	_, err = tx.Exec(ddlCreateSurah)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ddlCreateWord)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(ddlCreateTracker)
	if err != nil {
		return nil, err
	}

	// Populate data
	_, err = tx.Exec(sqlInsertSurah)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(sqlInsertWord)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO tracker (id, last_word)
		VALUES (1, 0)
		ON CONFLICT DO NOTHING`)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	err = tx.Commit()
	return db, err
}
