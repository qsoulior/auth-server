package jwt

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Fingerprint string   `json:"fingerprint"`
	Roles       []string `json:"roles"`
	jwt.RegisteredClaims
}
