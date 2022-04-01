//go:build !dev

package cmd

import (
	ap "github.com/muesli/go-app-paths"
)

func getDBPath() (string, error) {
	// Get platform specific path
	scope := ap.NewScope(ap.User, "kalimah")

	// Set database path
	dbPath, err := scope.DataPath("kalimah.db")
	if err != nil {
		return "", err
	}

	return dbPath, nil
}
