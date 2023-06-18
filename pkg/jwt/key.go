package jwt

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type key interface {
	HMAC() (any, error)
	RSA() (any, error)
	ECDSA() (any, error)
	Ed25519() (any, error)
}

type publicKey struct {
	data []byte
}

func (p publicKey) HMAC() (any, error) {
	return p.data, nil
}

func (p publicKey) RSA() (any, error) {
	return jwt.ParseRSAPublicKeyFromPEM(p.data)
}

func (p publicKey) ECDSA() (any, error) {
	return jwt.ParseECPublicKeyFromPEM(p.data)
}

func (p publicKey) Ed25519() (any, error) {
	return jwt.ParseEdPublicKeyFromPEM(p.data)
}

type privateKey struct {
	data []byte
}

func (p privateKey) HMAC() (any, error) {
	return p.data, nil
}

func (p privateKey) RSA() (any, error) {
	return jwt.ParseRSAPrivateKeyFromPEM(p.data)
}

func (p privateKey) ECDSA() (any, error) {
	return jwt.ParseECPrivateKeyFromPEM(p.data)
}

func (p privateKey) Ed25519() (any, error) {
	return jwt.ParseEdPrivateKeyFromPEM(p.data)
}

type keyParser struct {
	key key
}

func (p keyParser) Parse(alg string) (any, error) {
	method, err := GetSigningMethod(alg)

	switch method.(type) {
	case *jwt.SigningMethodHMAC:
		return p.key.HMAC()
	case *jwt.SigningMethodRSA, *jwt.SigningMethodRSAPSS:
		return p.key.RSA()
	case *jwt.SigningMethodECDSA:
		return p.key.ECDSA()
	case *jwt.SigningMethodEd25519:
		return p.key.Ed25519()
	}

	return nil, err
}

type parseFunc func(data []byte, alg string) (any, error)

func ParsePublicKey(data []byte, alg string) (any, error) {
	parser := keyParser{publicKey{data}}
	return parser.Parse(alg)
}

func ParsePrivateKey(data []byte, alg string) (any, error) {
	parser := keyParser{privateKey{data}}
	return parser.Parse(alg)
}

type readFunc func(path string, alg string) (any, error)

func read(p parseFunc) readFunc {
	return func(path string, alg string) (any, error) {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return p(data, alg)
	}
}

func ReadPublicKey(path string, alg string) (any, error) {
	return read(ParsePublicKey)(path, alg)
}

func ReadPrivateKey(path string, alg string) (any, error) {
	return read(ParsePrivateKey)(path, alg)
}
