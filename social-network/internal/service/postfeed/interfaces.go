package postfeed

import (
	"context"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
)

type Cache interface {
	GetUserFeed(ctx context.Context, input dto.GetPostFeed) ([]domain.Post, bool, error)
	SetUserFeed(ctx context.Context, userID uuid.UUID, posts []domain.Post) error
	UserFeedExists(ctx context.Context, userID uuid.UUID) (bool, error)
}
