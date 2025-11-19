package handlers

import (
	"errors"

	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type Handlers struct {
	auth AuthService
}

func New(auth AuthService) (*Handlers, error) {
	if utils.IsNil(auth) {
		return nil, errors.New("auth service is nil")
	}

	return &Handlers{
		auth: auth,
	}, nil
}
