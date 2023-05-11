package app

import (
	"fmt"

	"github.com/qsoulior/auth-server/pkg/jwt"
)

func NewJWT(privatePath string, publicPath string, issuer string, algorithm string) (jwt.Builder, jwt.Parser, error) {
	privateKey, err := jwt.ReadPrivateKey(privatePath, algorithm)
	if err != nil {
		return nil, nil, fmt.Errorf("private key: %w", err)
	}

	publicKey, err := jwt.ReadPublicKey(publicPath, algorithm)
	if err != nil {
		return nil, nil, fmt.Errorf("public key: %w", err)
	}

	builder, err := jwt.NewBuilder(jwt.Params{issuer, algorithm, privateKey})
	if err != nil {
		return nil, nil, fmt.Errorf("builder: %w", err)
	}

	parser, err := jwt.NewParser(jwt.Params{issuer, algorithm, publicKey})
	if err != nil {
		return nil, nil, fmt.Errorf("parser: %w", err)
	}

	return builder, parser, nil
}
