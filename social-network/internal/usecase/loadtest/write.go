package loadtest

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type UseCases struct {
	repo domain.LoadTestRepository
}

func New(repo domain.LoadTestRepository) (*UseCases, error) {
	if utils.IsNil(repo) {
		return nil, errors.New("repository is nil")
	}

	return &UseCases{repo: repo}, nil
}

func (u *UseCases) Write(ctx context.Context, value string) error {
	err := u.repo.Insert(ctx, uuid.New(), value)
	if err != nil {
		return fmt.Errorf("failed to insert load test value: %w", err)
	}

	return nil
}
