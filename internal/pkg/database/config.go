package database

import "fmt"

// Config keeps config to work with database.
type Config struct {
	Host           string `env:"DB_HOST,required"`
	Port           int    `env:"DB_PORT,required"`
	Name           string `env:"DB_NAME,required"`
	User           string `env:"DB_USER,required"`
	Password       string `env:"DB_PASSWORD,required"`
	MaxConnections int    `env:"DB_MAX_CONN,required"`
	EnableLog      bool   `env:"DB_ENABLE_LOG"`
}

func (c *Config) addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
