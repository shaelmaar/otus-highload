package post

import (
	"context"
	"errors"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type TxExecutor interface {
	Exec(
		ctx context.Context,
		f func(ctx context.Context, tx transaction.Tx) error,
		rollbackFn func(ctx context.Context),
	) error
}

type UseCases struct {
	repo domain.PostRepository
	tx   TxExecutor
}

func New(repo domain.PostRepository, tx TxExecutor) (*UseCases, error) {
	if utils.IsNil(repo) {
		return nil, errors.New("repo is nil")
	}

	if utils.IsNil(tx) {
		return nil, errors.New("tx is nil")
	}

	return &UseCases{
		repo: repo,
		tx:   tx,
	}, nil
}
