//go:build !prod

package cmd

func init() {
	developmentMode = true
}

func getDBPath() (string, error) {
	return "kalimah.db", nil
}
