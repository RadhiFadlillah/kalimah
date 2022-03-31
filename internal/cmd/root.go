package cmd

import (
	"fmt"
	"kalimah/internal/database"
	"os"
	fp "path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

var (
	db     *sqlx.DB
	dbPath string
)

// RootCmd returns the root command for yla-client
func RootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "kalimah",
		Short: "Simple web app for memorizing Quran translation",

		RunE:               runHandler,
		PersistentPreRunE:  preRunHandler,
		PersistentPostRunE: postRunHandler,
	}
}

func preRunHandler(cmd *cobra.Command, args []string) error {
	// Get and create database path
	dbPath, err := getDBPath()
	if err != nil {
		return fmt.Errorf("failed to get database path: %w", err)
	}

	err = os.MkdirAll(fp.Dir(dbPath), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create database dir: %w", err)
	}

	// Open database
	db, err = database.Open(dbPath)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	return nil
}

func postRunHandler(cmd *cobra.Command, args []string) error {
	return db.Close()
}

func runHandler(cmd *cobra.Command, args []string) error {
	return nil
}
