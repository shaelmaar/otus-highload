package config

import (
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

type Auth struct {
	SecretKey  string        `envconfig:"SECRET_KEY" required:"true"`
	Expiration time.Duration `envconfig:"EXPIRATION" default:"30m"`
}

func (a Auth) Validate() error {
	if len(a.SecretKey) < 32 {
		return errors.New("secret key must be at least 32 bytes long")
	}

	decodedKey, err := base64.URLEncoding.DecodeString(a.SecretKey)

	switch {
	case err != nil:
		return fmt.Errorf("failed to decode base64 secret key: %w", err)
	case len(decodedKey) != 32:
		return errors.New("secret key must be at least 32 bytes long")
	}

	return nil
}
