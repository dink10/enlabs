package main

import (
	"github.com/sirupsen/logrus"

	"github.com/dink10/enlabs/info"
	"github.com/dink10/enlabs/internal/app/api"
)

func main() {
	logrus.Infof("Application version: %s", info.Version)

	if err := api.Run(); err != nil {
		logrus.Fatal(err)
	}
}
