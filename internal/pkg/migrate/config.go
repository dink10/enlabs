package migrate

import (
	"github.com/dink10/enlabs/internal/pkg/database"
	"github.com/dink10/enlabs/internal/pkg/logger"
)

// Config is an application config.
type Config struct {
	Logger   logger.Config
	Database database.Config
}
