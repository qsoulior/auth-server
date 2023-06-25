package app

import (
	"fmt"

	"github.com/qsoulior/auth-server/pkg/jwt"
)

// NewJWT reads private and public keys, creates JWT builder and JWT parser.
// It returns error if keys read failed or configuration is incorrect.
func NewJWT(cfg *Config) (jwt.Builder, jwt.Parser, error) {
	privateKey, err := jwt.ReadPrivateKey(cfg.Key.PrivatePath, cfg.AT.Alg)
	if err != nil {
		return nil, nil, fmt.Errorf("private key: %w", err)
	}

	publicKey, err := jwt.ReadPublicKey(cfg.Key.PublicPath, cfg.AT.Alg)
	if err != nil {
		return nil, nil, fmt.Errorf("public key: %w", err)
	}

	builder, err := jwt.NewBuilder(jwt.Params{cfg.Name, cfg.AT.Alg, privateKey})
	if err != nil {
		return nil, nil, fmt.Errorf("builder: %w", err)
	}

	parser, err := jwt.NewParser(jwt.Params{cfg.Name, cfg.AT.Alg, publicKey})
	if err != nil {
		return nil, nil, fmt.Errorf("parser: %w", err)
	}

	return builder, parser, nil
}
