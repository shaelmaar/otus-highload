package transaction

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type Tx interface {
	pgx.Tx
}

type TxExecutor struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) (*TxExecutor, error) {
	if utils.IsNil(db) {
		return nil, errors.New("db is nil")
	}

	return &TxExecutor{
		db: db,
	}, nil
}

// Exec запускает функцию f в новой транзакции, если не передана явно.
func (t *TxExecutor) Exec(
	ctx context.Context,
	f func(ctx context.Context, tx Tx) error,
	rollbackFn func(ctx context.Context),
) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	err = f(ctx, tx)
	if err != nil {
		_ = tx.Rollback(ctx)

		if rollbackFn != nil {
			rollbackFn(ctx)
		}

		return fmt.Errorf("failed to execute transaction: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
