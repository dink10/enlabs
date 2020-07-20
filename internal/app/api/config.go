package api

import (
	"github.com/dink10/enlabs/internal/pkg/database"
	"github.com/dink10/enlabs/internal/pkg/logger"
	"github.com/dink10/enlabs/internal/pkg/server"
)

// Config is an application config.
type Config struct {
	Server   server.Config
	Logger   logger.Config
	Database database.Config
}
