package main

import (
	"kalimah/internal/cmd"

	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

func main() {
	// Format logrus
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	// Execute cmd
	err := cmd.RootCmd(assets).Execute()
	if err != nil {
		logrus.Fatalln(err)
	}
}
