package app

import (
	"github.com/qsoulior/auth-server/pkg/config"
)

const (
	EnvDev  Environment = "development"
	EnvProd Environment = "production"
)

type (
	// Config represents app configuration structure.
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
		PrivatePath string `env:"KEY_PRIVATE"`
		PublicPath  string `env:"KEY_PUBLIC"`
	}

	HTTPConfig struct {
		Host           string   `env:"HTTP_HOST" default:"0.0.0.0"`
		Port           string   `env:"HTTP_PORT" default:"3000"`
		AllowedOrigins []string `env:"HTTP_ORIGINS" default:"*"`
	}

	PostgresConfig struct {
		URI string `env:"POSTGRES_URI"`
	}

	ATConfig struct {
		Alg string `env:"AT_ALG" default:"HS256"`
		Age int    `env:"AT_AGE" default:"15"`
	}

	RTConfig struct {
		Cap int `env:"RT_CAP" default:"10"`
		Age int `env:"RT_AGE" default:"30"`
	}

	BcryptConfig struct {
		Cost int `env:"BCRYPT_COST" default:"4"`
	}
)

// NewConfig reads variables from file or environment
// and sets Config fields to variables.
// It returns pointer to a Config instance or nil if read failed.
func NewConfig(path string) (*Config, error) {
	cfg := new(Config)
	if path == "" {
		if err := config.ReadEnv(cfg); err != nil {
			return nil, err
		}
	} else {
		if err := config.ReadEnvFile(cfg, path); err != nil {
			return nil, err
		}
	}
	return cfg, nil
}
