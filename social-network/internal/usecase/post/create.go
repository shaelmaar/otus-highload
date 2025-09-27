package post

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
)

func (u *UseCases) Create(ctx context.Context, input dto.PostCreate) (uuid.UUID, error) {
	post := domain.Post{
		ID:           uuid.New(),
		Content:      input.Content,
		AuthorUserID: input.UserID,
		CreatedAt:    time.Now(),

		// не используется.
		UpdatedAt: time.Now(),
	}

	err := u.repo.Create(ctx, post)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to create post: %w", err)
	}

	return post.ID, nil
}
