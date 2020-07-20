package config

import (
	"github.com/caarlos0/env/v6"
)

// LoadConfig load config from environment variables.
func LoadConfig(cfg interface{}) error {
	return env.Parse(cfg)
}
