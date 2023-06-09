package jwt

import "github.com/golang-jwt/jwt/v5"

type Claims struct {
	Fingerprint string `json:"fingerprint"`
	jwt.RegisteredClaims
}
