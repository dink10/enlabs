package database

import (
	"context"
	"fmt"
	"time"

	"github.com/go-pg/pg/v9"
	"github.com/rafaeljesus/retry-go"
	"github.com/sirupsen/logrus"

	"github.com/dink10/enlabs/internal/pkg/logger"
)

const (
	retryInterval   = 500 * time.Millisecond
	numberOfRetries = 5
)

// Connect establishes connection to the database and pings it to ensure
// that connection was successful.
func Connect(ctx context.Context, cfg *Config) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     cfg.addr(),
		User:     cfg.User,
		Password: cfg.Password,
		Database: cfg.Name,
		PoolSize: cfg.MaxConnections,
	})

	if cfg.EnableLog {
		queryLogger := logger.NewDatabaseLogger()
		db.AddQueryHook(queryLogger)
	}

	err := retry.Do(pingFunc(ctx, db), numberOfRetries, retryInterval)
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

// Close is a helper function to close database connection and logrus error
// in case of failure.
func Close(db *pg.DB) {
	if err := db.Close(); err != nil {
		logrus.Errorf("failed to close database: %v", err)
	}
}

func pingFunc(ctx context.Context, db *pg.DB) retry.Func {
	const pingQuery = "SELECT 1"
	const pingTimeout = 500 * time.Millisecond

	return func() error {
		pingCtx, cancelPing := context.WithTimeout(ctx, pingTimeout)
		defer cancelPing()

		_, err := db.ExecContext(pingCtx, pingQuery)
		return err
	}
}
