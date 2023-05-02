package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Builder interface {
	Build(subject string) (string, error)
}

type builder struct {
	*handler
}

func NewBuilder(issuer string, algorithm string, key any) (*builder, error) {
	return new[builder](issuer, algorithm, key)
}

func (b *builder) Build(subject string) (string, error) {
	method := jwt.GetSigningMethod(b.algorithm)
	claims := &jwt.RegisteredClaims{
		Issuer:    b.issuer,
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
	}

	token := jwt.NewWithClaims(method, claims)
	return token.SignedString(b.key)
}
