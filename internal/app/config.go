package app

import (
	"fmt"

	"github.com/qsoulior/auth-server/pkg/config"
)

const (
	EnvDev  Environment = "development"
	EnvProd Environment = "production"
)

type (
	Environment string
	Config      struct {
		Env      Environment `env:"APP_ENV" default:"development"`
		HTTP     HTTPConfig
		Postgres PostgresConfig
	}

	HTTPConfig struct {
		Port string `env:"HTTP_PORT" default:"8000"`
	}

	PostgresConfig struct {
		URI string `env:"POSTGRES_URI"`
	}
)

func NewConfig() (*Config, error) {
	cfg := new(Config)
	if err := config.ReadEnvFile(cfg, "configs/dev.env"); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	return cfg, nil
}
