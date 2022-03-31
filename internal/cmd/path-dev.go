//go:build !prod

package cmd

func getDBPath() (string, error) {
	return "kalimah.db", nil
}
