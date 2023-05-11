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
	params Params
	method jwt.SigningMethod
}

func NewParser(params Params) (*parser, error) {
	method, err := GetSigningMethod(params.Algorithm)
	if err != nil {
		return nil, err
	}

	return &parser{params, method}, nil
}

func (p *parser) Parse(tokenStr string) (string, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{p.params.Algorithm}), jwt.WithIssuer(p.params.Issuer))
	token, err := parser.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		return p.params.Key, nil
	})

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		return claims.Subject, nil
	}

	return "", ErrClaimsInvalid
}
