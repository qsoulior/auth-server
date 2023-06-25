package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Builder is interface implemented by types that can build JWT.
type Builder interface {
	// Build creates new Claims, JWT string and signing it.
	// It returns JWT string or empty string if signing failed.
	Build(subject string, age time.Duration, fingerprint string, roles []string) (string, error)
}

// builder implements Builder interface.
type builder struct {
	params Params
	method jwt.SigningMethod
}

// NewBuilder gets signing method by params and creates new builder.
// It returns pointer to a builder instance or nil if params.Algorithm is incorrect.
func NewBuilder(params Params) (*builder, error) {
	method, err := GetSigningMethod(params.Algorithm)
	if err != nil {
		return nil, err
	}

	return &builder{params, method}, nil
}

// Build creates new Claims, JWT string and signing it.
// It returns JWT string or empty string if signing failed.
func (b *builder) Build(subject string, age time.Duration, fingerprint string, roles []string) (string, error) {
	claims := &Claims{
		fingerprint,
		roles,
		jwt.RegisteredClaims{
			Issuer:    b.params.Issuer,
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(age)),
		},
	}

	token := jwt.NewWithClaims(b.method, claims)
	return token.SignedString(b.params.Key)
}
