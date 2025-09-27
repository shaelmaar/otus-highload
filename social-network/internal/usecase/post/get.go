package post

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
)

func (u *UseCases) GetByID(ctx context.Context, id uuid.UUID) (domain.Post, error) {
	var out domain.Post

	post, err := u.repo.Slave().GetByID(ctx, id)
	if err != nil {
		return out, fmt.Errorf("failed to get post by id: %w", err)
	}

	return post, nil
}
