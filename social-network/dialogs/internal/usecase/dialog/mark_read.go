package dialog

import (
	"context"
	"fmt"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/domain"
)

func (u *UseCases) MarkMessagesAsRead(ctx context.Context, dialogID string, messageIDs []uint64) error {
	_, err := u.repo.UpdateMessagesStateFrom(
		ctx, dialogID, messageIDs, domain.DialogMessageStateReading, domain.DialogMessageStateRead)
	if err != nil {
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}

	return nil
}
