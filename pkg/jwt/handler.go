package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

var ErrAlgInvalid = errors.New("algorithm is invalid")

type handler struct {
	issuer    string
	algorithm string
	key       any
}

func newHandler(issuer string, algorithm string, key any) (*handler, error) {
	for _, alg := range jwt.GetAlgorithms() {
		if algorithm == alg {
			return &handler{issuer, algorithm, key}, nil
		}
	}
	return nil, ErrAlgInvalid
}

func new[T builder | parser](issuer string, algorithm string, key any) (*T, error) {
	handler, err := newHandler(issuer, algorithm, key)
	if err != nil {
		return nil, err
	}
	return &T{handler}, nil
}
