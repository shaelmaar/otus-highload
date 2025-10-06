package feed

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (u *UseCases) UpdateUserFeedChunked(ctx context.Context, userIDs []uuid.UUID) error {
	for _, userID := range userIDs {
		err := u.service.UpdateUserFeedCache(ctx, userID)
		if err != nil {
			return fmt.Errorf("failed to update user feed cache: %w", err)
		}
	}

	return nil
}
