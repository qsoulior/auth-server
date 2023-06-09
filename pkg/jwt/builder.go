package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Builder interface {
	Build(subject string, age time.Duration, fingerprint string) (string, error)
}

type builder struct {
	params Params
	method jwt.SigningMethod
}

func NewBuilder(params Params) (*builder, error) {
	method, err := GetSigningMethod(params.Algorithm)
	if err != nil {
		return nil, err
	}

	return &builder{params, method}, nil
}

func (b *builder) Build(subject string, age time.Duration, fingerprint string) (string, error) {
	claims := &Claims{
		fingerprint,
		jwt.RegisteredClaims{
			Issuer:    b.params.Issuer,
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(age)),
		},
	}

	token := jwt.NewWithClaims(b.method, claims)
	return token.SignedString(b.params.Key)
}
