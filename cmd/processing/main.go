package main

import (
	"github.com/sirupsen/logrus"

	"github.com/dink10/enlabs/internal/app/processing"
)

// revision of application
var revision string

func main() {
	logrus.Infof("Processing revision: %s", revision)

	if err := processing.Run(); err != nil {
		logrus.Fatal(err)
	}
}
