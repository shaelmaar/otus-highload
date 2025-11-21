package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTConfig struct {
	SecretKey  string
	Expiration time.Duration
	Issuer     string
}
