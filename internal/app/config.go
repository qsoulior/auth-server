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
	Config struct {
		Name     string      `env:"APP_NAME" default:"auth"`
		Env      Environment `env:"APP_ENV" default:"development"`
		HTTP     HTTPConfig
		JWT      JWTConfig
		Bcrypt   BcryptConfig
		Postgres PostgresConfig
	}

	Environment string

	HTTPConfig struct {
		Port string `env:"HTTP_PORT" default:"8000"`
	}

	JWTConfig struct {
		Alg string `env:"JWT_ALG" default:"HS256"`
	}

	BcryptConfig struct {
		Cost int `env:"BCRYPT_COST" default:"4"`
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
