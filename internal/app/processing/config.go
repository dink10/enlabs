package processing

import (
	"github.com/dink10/enlabs/internal/pkg/database"
	"github.com/dink10/enlabs/internal/pkg/logger"
	"github.com/dink10/enlabs/internal/pkg/processing"
)

// Config is an application config.
type Config struct {
	Processing processing.Config
	Logger     logger.Config
	Database   database.Config
}
