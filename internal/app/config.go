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
		Key      KeyConfig
		HTTP     HTTPConfig
		Postgres PostgresConfig
		AT       ATConfig
		RT       RTConfig
		Bcrypt   BcryptConfig
	}

	Environment string

	KeyConfig struct {
		PrivatePath string `env:"KEY_PRIVATE_PATH"`
		PublicPath  string `env:"KEY_PUBLIC_PATH"`
	}

	HTTPConfig struct {
		Port string `env:"HTTP_PORT" default:"8000"`
	}

	PostgresConfig struct {
		URI string `env:"POSTGRES_URI"`
	}

	RTConfig struct {
		Cap int `env:"RT_CAP" default:"10"`
		Age int `env:"RT_AGE" default:"30"`
	}

	ATConfig struct {
		Alg string `env:"AT_ALG" default:"HS256"`
		Age int    `env:"AT_AGE" default:"15"`
	}

	BcryptConfig struct {
		Cost int `env:"BCRYPT_COST" default:"4"`
	}
)

func NewConfig() (*Config, error) {
	cfg := new(Config)
	if err := config.ReadEnvFile(cfg, "configs/dev.env"); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	return cfg, nil
}
