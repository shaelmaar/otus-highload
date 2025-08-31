package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
)

func (uc *UseCases) GetByID(ctx context.Context, id uuid.UUID) (domain.User, error) {
	var (
		out domain.User
		err error
	)

	out, err = uc.repo.Slave().GetByID(ctx, id)
	if err != nil {
		return out, fmt.Errorf("failed to get user by id: %w", err)
	}

	return out, nil
}
