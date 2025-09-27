package user

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
)

type TxExecutor interface {
	Exec(
		ctx context.Context,
		f func(ctx context.Context, tx transaction.Tx) error,
		rollbackFn func(ctx context.Context),
	) error
}

type AuthService interface {
	GenerateToken(userID string) (string, error)
}
