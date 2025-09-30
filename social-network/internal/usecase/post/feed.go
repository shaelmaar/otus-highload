package post

import (
	"context"
	"fmt"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
)

func (u *UseCases) GetPostFeed(ctx context.Context, input dto.GetPostFeed) ([]domain.Post, error) {
	friendIDs, err := u.friendRepo.Slave().GetUserFriendIDs(ctx, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user friend ids: %w", err)
	}

	posts, err := u.repo.Slave().GetLastPostsByUserIDs(ctx, dto.GetLastPostsByUserIDs{
		UserIDs: friendIDs,
		Offset:  input.Offset,
		Limit:   input.Limit,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get last posts by user ids: %w", err)
	}

	return posts, nil
}
