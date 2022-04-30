package cmd

import (
	"fmt"
	"kalimah/internal/backend"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func startCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start the server",
		RunE:  startCmdHandler,
	}

	cmd.Flags().IntP("port", "p", 8080, "Port used by the server")
	return cmd
}

func startCmdHandler(cmd *cobra.Command, args []string) error {
	// Get flags value
	port, _ := cmd.Flags().GetInt("port")

	// Start server
	server := backend.Server{
		DB:      db,
		Assets:  assets,
		DevMode: developmentMode,
	}

	if developmentMode {
		logrus.Println("development mode enabled")
	}

	if err := server.Serve(port); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
