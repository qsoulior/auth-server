package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var ErrClaimsInvalid = errors.New("claims is invalid")

type Parser interface {
	Parse(tokenString string) (*Claims, error)
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

func (p *parser) Parse(tokenString string) (*Claims, error) {
	parser := jwt.NewParser(jwt.WithValidMethods([]string{p.params.Algorithm}), jwt.WithIssuer(p.params.Issuer))
	token, err := parser.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (any, error) {
		return p.params.Key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok {
		return claims, nil
	}

	return nil, ErrClaimsInvalid
}
