package logger

import (
	"context"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/sirupsen/logrus"
)

// NewDatabaseLogger returns database logger.
func NewDatabaseLogger() *DatabaseLogger {
	return &DatabaseLogger{
		Logger: logrus.StandardLogger(),
	}
}

// DatabaseLogger implements pg.QueryHook interface and is used for logging database queries.
type DatabaseLogger struct {
	*logrus.Logger
}

// BeforeQuery basically skips this step doing nothing.
func (l *DatabaseLogger) BeforeQuery(ctx context.Context, _ *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

// AfterQuery logs executed query.
func (l *DatabaseLogger) AfterQuery(_ context.Context, q *pg.QueryEvent) error {
	query, err := q.UnformattedQuery()
	if err != nil {
		return err
	}

	query = strings.ReplaceAll(query, "\"", "'")
	l.WithField("query", query).Debug("query executed")

	return nil
}
