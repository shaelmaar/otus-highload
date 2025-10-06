package post

import (
	"context"
	"fmt"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
)

func (u *UseCases) GetPostFeed(ctx context.Context, input dto.GetPostFeed) ([]domain.Post, error) {
	posts, err := u.feedService.GetUserFeed(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get user feed: %w", err)
	}

	return posts, nil
}
