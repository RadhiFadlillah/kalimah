package main

import (
	"kalimah/internal/cmd"

	"github.com/sirupsen/logrus"
	_ "modernc.org/sqlite"
)

func main() {
	// Format logrus
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Execute cmd
	err := cmd.RootCmd().Execute()
	if err != nil {
		logrus.Fatalln(err)
	}
}
