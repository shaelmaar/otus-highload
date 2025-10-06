package feed

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (u *UseCases) UpdateUserFeed(ctx context.Context, userID uuid.UUID) error {
	err := u.service.UpdateUserFeedCache(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to update user feed cache: %w", err)
	}

	return nil
}
