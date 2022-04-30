package cmd

import (
	"fmt"
	"io/fs"
	"kalimah/internal/database"
	"os"
	fp "path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

var (
	db     *sqlx.DB
	assets fs.FS

	developmentMode = false
)

// RootCmd returns the root command for yla-client
func RootCmd(assetsFs fs.FS) *cobra.Command {
	assets = assetsFs
	rootCmd := &cobra.Command{
		Use:   "kalimah",
		Short: "Simple web app for memorizing Quran translation",

		PersistentPreRunE:  preRunHandler,
		PersistentPostRunE: postRunHandler,
	}

	rootCmd.AddCommand(startCmd(), initCmd())
	return rootCmd
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
