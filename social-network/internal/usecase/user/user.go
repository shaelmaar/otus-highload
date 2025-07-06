package user

import (
	"errors"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type UseCases struct {
	repo domain.UserRepository
	tx   TxExecutor
}

func New(
	repo domain.UserRepository,
	tx TxExecutor,
) (*UseCases, error) {
	if utils.IsNil(repo) {
		return nil, errors.New("user repository is nil")
	}

	if utils.IsNil(tx) {
		return nil, errors.New("tx executor is nil")
	}

	return &UseCases{
		repo: repo,
		tx:   tx,
	}, nil
}
