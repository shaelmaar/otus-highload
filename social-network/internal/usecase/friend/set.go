package friend

import (
	"context"
	"fmt"
	"time"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
)

func (u *UseCases) Set(ctx context.Context, input dto.FriendSet) error {
	err := u.repo.Create(ctx, domain.Friend{
		UserID:    input.UserID,
		FriendID:  input.FriendID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to create friend %w", err)
	}

	return nil
}
