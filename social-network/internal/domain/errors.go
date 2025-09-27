package domain

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrPostNotFound        = errors.New("post not found")
	ErrPostNotFoundForUser = errors.New("post not found for user")
)
