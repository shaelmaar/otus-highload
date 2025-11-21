package user

import (
	"errors"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type UseCases struct {
	repo domain.UserRepository
	auth AuthService
	tx   TxExecutor
}

func New(
	repo domain.UserRepository,
	auth AuthService,
	tx TxExecutor,
) (*UseCases, error) {
	if utils.IsNil(repo) {
		return nil, errors.New("user repository is nil")
	}

	if utils.IsNil(auth) {
		return nil, errors.New("auth service is nil")
	}

	if utils.IsNil(tx) {
		return nil, errors.New("tx executor is nil")
	}

	return &UseCases{
		repo: repo,
		auth: auth,
		tx:   tx,
	}, nil
}
