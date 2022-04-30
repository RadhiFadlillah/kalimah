package cmd

import (
	"kalimah/internal/database"

	"github.com/spf13/cobra"
)

func initCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initiate the database",
		RunE:  initCmdHandler,
	}
}

func initCmdHandler(cmd *cobra.Command, args []string) error {
	return database.PopulateData(db)
}
