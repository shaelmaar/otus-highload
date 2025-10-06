package feed

import (
	"errors"

	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type UseCases struct {
	service Service
}

func New(service Service) (*UseCases, error) {
	if utils.IsNil(service) {
		return nil, errors.New("service is nil")
	}

	return &UseCases{
		service: service,
	}, nil
}
