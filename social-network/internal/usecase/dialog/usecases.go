package dialog

import (
	"errors"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type UseCases struct {
	repo domain.DialogRepository
}

func New(repo domain.DialogRepository) (*UseCases, error) {
	if utils.IsNil(repo) {
		return nil, errors.New("repository is nil")
	}

	return &UseCases{
		repo: repo,
	}, nil
}
