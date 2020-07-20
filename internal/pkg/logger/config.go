package logger

// Config keeps configuration of logger.
type Config struct {
	Level string `env:"LOG_LEVEL,required"`
}
