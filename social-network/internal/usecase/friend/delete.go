package friend

import (
	"context"
	"fmt"
	"time"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
)

func (u *UseCases) Delete(ctx context.Context, input dto.FriendDelete) error {
	err := u.repo.Delete(ctx, domain.Friend{
		UserID:   input.UserID,
		FriendID: input.FriendID,

		// не используется.
		CreatedAt: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("failed to delete friend %w", err)
	}

	return nil
}
