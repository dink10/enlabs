package main

import (
	"context"
	"os"

	"github.com/dink10/enlabs/internal/pkg/config"
	"github.com/dink10/enlabs/internal/pkg/database"
	"github.com/dink10/enlabs/internal/pkg/migrate"
	migrations "github.com/robinjoseph08/go-pg-migrations/v2"
	"github.com/sirupsen/logrus"
)

const directory = "tools/migrations"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var cfg migrate.Config
	err := config.LoadConfig(&cfg)
	if err != nil {
		logrus.Fatalf("failed to parse config: %v", err)
	}

	db, err := database.Connect(ctx, &cfg.Database)
	if err != nil {
		logrus.Fatalf("failed to connect to database: %v", err)
	}
	defer database.Close(db)

	err = migrations.Run(db, directory, os.Args)
	if err != nil {
		logrus.Fatal(err)
	}

	return
}
