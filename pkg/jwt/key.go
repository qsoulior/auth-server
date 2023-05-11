package jwt

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

func ReadPublicKey(path string, alg string) (any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParsePublicKey(data, alg)
}

func ReadPrivateKey(path string, alg string) (any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParsePrivateKey(data, alg)
}

func ParsePublicKey(data []byte, alg string) (any, error) {
	parser := &publicParser{data}
	return parseKey(parser, alg)
}

func ParsePrivateKey(data []byte, alg string) (any, error) {
	parser := &privateParser{data}
	return parseKey(parser, alg)
}

type keyParser interface {
	HMAC() (any, error)
	RSA() (any, error)
	ECDSA() (any, error)
	Ed25519() (any, error)
}

type publicParser struct {
	data []byte
}

func (p *publicParser) HMAC() (any, error) {
	return p.data, nil
}

func (p *publicParser) RSA() (any, error) {
	return jwt.ParseRSAPublicKeyFromPEM(p.data)
}

func (p *publicParser) ECDSA() (any, error) {
	return jwt.ParseECPublicKeyFromPEM(p.data)
}

func (p *publicParser) Ed25519() (any, error) {
	return jwt.ParseEdPublicKeyFromPEM(p.data)
}

type privateParser struct {
	data []byte
}

func (p *privateParser) HMAC() (any, error) {
	return p.data, nil
}

func (p *privateParser) RSA() (any, error) {
	return jwt.ParseRSAPrivateKeyFromPEM(p.data)
}

func (p *privateParser) ECDSA() (any, error) {
	return jwt.ParseECPrivateKeyFromPEM(p.data)
}

func (p *privateParser) Ed25519() (any, error) {
	return jwt.ParseEdPrivateKeyFromPEM(p.data)
}

func parseKey(parser keyParser, alg string) (any, error) {
	method, err := GetSigningMethod(alg)
	if err != nil {
		return nil, err
	}

	switch method.(type) {
	case *jwt.SigningMethodHMAC:
		return parser.HMAC()
	case *jwt.SigningMethodRSA, *jwt.SigningMethodRSAPSS:
		return parser.RSA()
	case *jwt.SigningMethodECDSA:
		return parser.ECDSA()
	case *jwt.SigningMethodEd25519:
		return parser.Ed25519()
	}

	return nil, ErrAlgInvalid
}
