package app

import (
	"fmt"

	"github.com/qsoulior/auth-server/pkg/jwt"
)

func NewJWT(cfg *Config) (jwt.Builder, jwt.Parser, error) {
	privateKey, err := jwt.ReadPrivateKey(cfg.Key.PrivatePath, cfg.JWT.Alg)
	if err != nil {
		return nil, nil, fmt.Errorf("private key: %w", err)
	}

	publicKey, err := jwt.ReadPublicKey(cfg.Key.PublicPath, cfg.JWT.Alg)
	if err != nil {
		return nil, nil, fmt.Errorf("public key: %w", err)
	}

	builder, err := jwt.NewBuilder(jwt.Params{cfg.Name, cfg.JWT.Alg, privateKey})
	if err != nil {
		return nil, nil, fmt.Errorf("builder: %w", err)
	}

	parser, err := jwt.NewParser(jwt.Params{cfg.Name, cfg.JWT.Alg, publicKey})
	if err != nil {
		return nil, nil, fmt.Errorf("parser: %w", err)
	}

	return builder, parser, nil
}
