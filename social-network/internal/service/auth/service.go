package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	jwtConfig JWTConfig
}

func NewService(
	secretKey string,
	expiration time.Duration,
	issuer string,
) (*Service, error) {
	if secretKey == "" {
		return nil, errors.New("secret key is empty")
	}

	if expiration == 0 {
		return nil, errors.New("expiration is zero")
	}

	if issuer == "" {
		return nil, errors.New("issuer is empty")
	}

	return &Service{
		jwtConfig: JWTConfig{
			SecretKey:  secretKey,
			Expiration: expiration,
			Issuer:     issuer,
		}}, nil
}

func (s *Service) GenerateToken(userID string) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		//nolint:exhaustruct // остальное не нужно.
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.jwtConfig.Expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.jwtConfig.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.jwtConfig.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (s *Service) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, new(CustomClaims), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.jwtConfig.SecretKey), nil
	})

	switch {
	case errors.Is(err, jwt.ErrTokenExpired):
		return "", ErrTokenExpired
	case err != nil:
		return "", fmt.Errorf("%s: %w", err.Error(), ErrTokenInvalid)
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return "", ErrTokenInvalid
}
