package jwt

import (
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// key is interface implemented by types
// that can represent keys encoded in PEM.
type key interface {
	// HMAC parses data and returns HMAC bytes.
	HMAC() (any, error)

	// RSA parses data and returns pointer to a rsa.PublicKey/rsa.PrivateKey.
	// It returns nil if key isn't RSA.
	RSA() (any, error)

	// ECDSA parses data and returns pointer to a ecdsa.PublicKey/ecdsa.PrivateKey.
	// It returns nil if key isn't ECDSA or key's BitSize isn't equal bitSize.
	ECDSA(bitSize int) (any, error)

	// Ed25519 parses data and returns crypto.PublicKey/crypto.PrivateKey.
	// It returns nil if key isn't Ed25519.
	Ed25519() (any, error)
}

// publicKey implements key interface.
// It represents public keys encoded in PEM.
type publicKey struct {
	data []byte
}

// HMAC parses p.data and returns HMAC bytes.
func (p publicKey) HMAC() (any, error) {
	return p.data, nil
}

// RSA parses p.data and returns pointer to a rsa.PublicKey.
// It returns nil if key isn't RSA.
func (p publicKey) RSA() (any, error) {
	return jwt.ParseRSAPublicKeyFromPEM(p.data)
}

// ECDSA parses p.data and returns pointer to a ecdsa.PublicKey.
// It returns nil if key isn't ECDSA or key's BitSize isn't equal bitSize.
func (p publicKey) ECDSA(bitSize int) (any, error) {
	key, err := jwt.ParseECPublicKeyFromPEM(p.data)
	if err != nil {
		return nil, err
	}
	if key.Params().BitSize != bitSize {
		return nil, jwt.ErrInvalidKey
	}
	return key, nil
}

// Ed25519 parses p.data and returns crypto.PublicKey.
// It returns nil if key isn't Ed25519.
func (p publicKey) Ed25519() (any, error) {
	return jwt.ParseEdPublicKeyFromPEM(p.data)
}

// privateKey implements key interface.
// It represents private keys encoded in PEM.
type privateKey struct {
	data []byte
}

// HMAC parses p.data and returns HMAC bytes.
func (p privateKey) HMAC() (any, error) {
	return p.data, nil
}

// RSA parses p.data and returns pointer to a rsa.PrivateKey.
// It returns nil if key isn't RSA.
func (p privateKey) RSA() (any, error) {
	return jwt.ParseRSAPrivateKeyFromPEM(p.data)
}

// ECDSA parses p.data and returns pointer to a ecdsa.PrivateKey.
// It returns nil if key isn't ECDSA or key's BitSize isn't equal bitSize.
func (p privateKey) ECDSA(bitSize int) (any, error) {
	key, err := jwt.ParseECPrivateKeyFromPEM(p.data)
	if err != nil {
		return nil, err
	}
	if key.Params().BitSize != bitSize {
		return nil, jwt.ErrInvalidKey
	}

	return key, nil
}

// Ed25519 parses p.data and returns crypto.PrivateKey.
// It returns nil if key isn't Ed25519.
func (p privateKey) Ed25519() (any, error) {
	return jwt.ParseEdPrivateKeyFromPEM(p.data)
}

// keyParser represents public/private key parser.
type keyParser struct {
	key key
}

// Parse gets jwt.SigningMethod by alg and calls one of the p.key methods.
// It returns key or nil if method is incorrect.
func (p keyParser) Parse(alg string) (any, error) {
	method, err := GetSigningMethod(alg)

	switch method := method.(type) {
	case *jwt.SigningMethodHMAC:
		return p.key.HMAC()
	case *jwt.SigningMethodRSA, *jwt.SigningMethodRSAPSS:
		return p.key.RSA()
	case *jwt.SigningMethodECDSA:
		return p.key.ECDSA(method.CurveBits)
	case *jwt.SigningMethodEd25519:
		return p.key.Ed25519()
	}

	return nil, err
}

// parseFunc represents function to parse public/private key.
type parseFunc func(data []byte, alg string) (any, error)

// ParsePublicKey creates keyParser and parses public key.
// It returns public key or nil if parsing failed.
func ParsePublicKey(data []byte, alg string) (any, error) {
	parser := keyParser{publicKey{data}}
	return parser.Parse(alg)
}

// ParsePrivateKey creates keyParser and parses private key.
// It returns private key or nil if parsing failed.
func ParsePrivateKey(data []byte, alg string) (any, error) {
	parser := keyParser{privateKey{data}}
	return parser.Parse(alg)
}

// readFunc represents function to read public/private key and parse it.
type readFunc func(path string, alg string) (any, error)

// read is parseFunc decorator.
// It returns readFunc.
func read(p parseFunc) readFunc {
	return func(path string, alg string) (any, error) {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		return p(data, alg)
	}
}

// ReadPublicKey reads public key and parses it.
// It returns public key or nil if reading/parsing failed
func ReadPublicKey(path string, alg string) (any, error) {
	return read(ParsePublicKey)(path, alg)
}

// ReadPrivateKey reads private key and parses it.
// It returns private key or nil if reading/parsing failed.
func ReadPrivateKey(path string, alg string) (any, error) {
	return read(ParsePrivateKey)(path, alg)
}
