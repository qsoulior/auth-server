package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrAlgInvalid = errors.New("algorithm is invalid")
	ErrAlgNone    = errors.New("algorithm is none")
)

type Params struct {
	Issuer    string
	Algorithm string
	Key       any
}

func GetSigningMethod(alg string) (jwt.SigningMethod, error) {
	method := jwt.GetSigningMethod(alg)

	if method == nil {
		return nil, ErrAlgInvalid
	}

	if method == jwt.SigningMethodNone {
		return nil, ErrAlgNone
	}

	return method, nil
}
