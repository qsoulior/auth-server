package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var ErrClaimsInvalid = errors.New("claims is invalid")

// Parser is interface implemented by types that can parse JWT.
type Parser interface {
	// Parse parses JWT string.
	// It returns pointer to a Claims or nil if parsing failed.
	Parse(tokenString string) (*Claims, error)
}

// parser implements Parser interface.
type parser struct {
	params Params
	method jwt.SigningMethod
}

// NewParser gets signing method by params and creates new parser.
// It returns pointer to a parser instance or nil if params.Algorithm is incorrect.
func NewParser(params Params) (*parser, error) {
	method, err := GetSigningMethod(params.Algorithm)
	if err != nil {
		return nil, err
	}

	return &parser{params, method}, nil
}

// Parse creates new jwt.Parser and parses JWT string.
// It returns pointer to a Claims or nil if parsing failed.
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
