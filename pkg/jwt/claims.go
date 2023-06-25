// Package jwt provides structures and functions to build/parse JWT,
// read/parse public and private keys.
package jwt

import "github.com/golang-jwt/jwt/v5"

// Claims represents custom claims.
// The jwt.RegisteredClaims embedded in it.
type Claims struct {
	Fingerprint string   `json:"fingerprint"`
	Roles       []string `json:"roles"`
	jwt.RegisteredClaims
}
