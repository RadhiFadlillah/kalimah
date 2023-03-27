package cmd

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/spf13/cobra"
)

var rxSurahAyah = regexp.MustCompile(`^(\d+):(\d+)`)

func markCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mark <surah>:<ayah>",
		Short: "Mark the last translated ayah",
		Args:  cobra.ExactArgs(1),
		RunE:  markCmdHandler,
	}
}

func markCmdHandler(cmd *cobra.Command, args []string) (err error) {
	// Parse args
	parts := rxSurahAyah.FindStringSubmatch(args[0])
	if len(parts) == 0 {
		return fmt.Errorf("argument %q is not in <surah>:<ayah> format", args[0])
	}

	surah, _ := strconv.ParseInt(parts[1], 10, 64)
	ayah, _ := strconv.ParseInt(parts[2], 10, 64)

	// Prepare transaction
	tx, err := db.Beginx()
	if err != nil {
		return
	}

	// Make sure to rollback if error ever happened
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Fetch the last word ID
	var lastWordID int
	err = tx.Get(&lastWordID,
		`WITH last_ayah AS (
			SELECT ?+start-1 ayah
			FROM surah WHERE id = ?)
		SELECT IFNULL(MAX(w.id), 0) last_word 
		FROM word w
		CROSS JOIN last_ayah la
		WHERE w.ayah = la.ayah`, ayah, surah)
	if err != nil {
		return
	}

	if lastWordID <= 0 {
		return fmt.Errorf("surah %d ayah %d not exist", surah, ayah)
	}

	// Save to track
	_, err = tx.Exec(
		`INSERT INTO tracker (id, last_word) VALUES (1, ?) 
		ON CONFLICT DO UPDATE SET last_word = excluded.last_word`,
		lastWordID)
	if err != nil {
		return err
	}

	// Commit
	return tx.Commit()
}
