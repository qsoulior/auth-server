package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var ErrClaimsInvalid = errors.New("claims is invalid")

type Parser interface {
	Parse(token string) (string, error)
}

type parser struct {
	*handler
}

func NewParser(issuer string, algorithm string, key any) (*parser, error) {
	return new[parser](issuer, algorithm, key)
}

func (p *parser) Parse(tokenStr string) (string, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{p.algorithm}), jwt.WithIssuer(p.issuer))

	token, err := parser.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return p.key, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		return claims.Subject, nil
	}
	return "", ErrClaimsInvalid
}
