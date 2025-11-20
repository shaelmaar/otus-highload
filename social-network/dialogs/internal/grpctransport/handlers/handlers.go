package handlers

import (
	"errors"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/pkg/utils"
)

type Handlers struct {
	dialogUseCases DialogUseCases
}

func New(dialogUseCases DialogUseCases) (*Handlers, error) {
	if utils.IsNil(dialogUseCases) {
		return nil, errors.New("dialog use cases is nil")
	}

	return &Handlers{
		dialogUseCases: dialogUseCases,
	}, nil
}
