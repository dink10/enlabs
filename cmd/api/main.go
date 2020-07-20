package main

import (
	"github.com/sirupsen/logrus"

	"github.com/dink10/enlabs/internal/app/api"
)

const (
	// version of application
	version = "1.0"
)

func main() {
	logrus.Infof("Application version: %s", version)

	if err := api.Run(); err != nil {
		logrus.Fatal(err)
	}
}
