package main

import (
	"github.com/sirupsen/logrus"

	"github.com/dink10/enlabs/info"
	"github.com/dink10/enlabs/internal/app/processing"
)

func main() {
	logrus.Infof("Processing version: %s", info.Version)

	if err := processing.Run(); err != nil {
		logrus.Fatal(err)
	}
}
