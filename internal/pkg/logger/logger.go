package logger

import (
	"github.com/sirupsen/logrus"
)

// Init initializes the application logger.
func Init(cfg *Config) error {
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		return err
	}

	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	return nil
}
