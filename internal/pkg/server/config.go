package server

import (
	"net"
)

// Config keeps configuration of HTTP server.
type Config struct {
	Host string `env:"SERVER_HOST,required"`
	Port string `env:"SERVER_PORT,required"`

	LogRequests bool `env:"SERVER_LOG_REQUESTS"`
}

func (c *Config) addr() string {
	return net.JoinHostPort(c.Host, c.Port)
}
