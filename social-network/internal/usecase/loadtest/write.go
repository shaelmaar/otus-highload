package loadtest

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type Metrics interface {
	IncLoadTestWrites(success bool)
}

type TxExecutor interface {
	Exec(
		ctx context.Context,
		f func(ctx context.Context, tx transaction.Tx) error,
		rollbackFn func(ctx context.Context),
	) error
}

type UseCases struct {
	repo       domain.LoadTestRepository
	txExecutor TxExecutor
	metrics    Metrics
	logger     *zap.Logger
}

func New(
	repo domain.LoadTestRepository,
	txExecutor TxExecutor,
	metrics Metrics,
	logger *zap.Logger,
) (*UseCases, error) {
	if utils.IsNil(repo) {
		return nil, errors.New("repository is nil")
	}

	if utils.IsNil(txExecutor) {
		return nil, errors.New("transaction executor is nil")
	}

	if utils.IsNil(metrics) {
		return nil, errors.New("metrics is nil")
	}

	if utils.IsNil(logger) {
		return nil, errors.New("logger is nil")
	}

	return &UseCases{
		repo:       repo,
		txExecutor: txExecutor,
		metrics:    metrics,
		logger:     logger,
	}, nil
}

func (u *UseCases) Write(ctx context.Context, value string) error {
	id := uuid.New()
	err := u.repo.Insert(ctx, id, value)

	if err != nil {
		//nolint:contextcheck // основной контекст может быть уже отменен.
		deleteErr := u.repo.Delete(context.Background(), id)
		if deleteErr != nil {
			u.logger.Warn("failed to delete loadtest",
				zap.Error(deleteErr), zap.String("id", id.String()))
		}

		u.metrics.IncLoadTestWrites(false)

		return fmt.Errorf("failed to insert load test value : %w", err)
	}

	u.metrics.IncLoadTestWrites(true)

	return nil
}
