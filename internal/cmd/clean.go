package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func cleanCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clean",
		Short: "Clean and reset database",
		RunE:  cleanCmdHandler,
	}
}

func cleanCmdHandler(cmd *cobra.Command, args []string) error {
	// Close database
	db.Close()

	// Remove files
	dbPath, err := getDBPath()
	if err != nil {
		return err
	}

	return os.RemoveAll(dbPath)
}
