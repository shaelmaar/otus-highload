package dialog

import (
	"context"
	"fmt"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/domain"
)

func (u *UseCases) MarkMessageAsFailed(ctx context.Context, dialogID string, messageID uint64) error {
	_, err := u.repo.UpdateMessageStateFrom(ctx, dialogID, messageID,
		domain.DialogMessageStateSending, domain.DialogMessageStateFailed)
	if err != nil {
		return fmt.Errorf("failed to mark message as sent: %w", err)
	}

	return nil
}
