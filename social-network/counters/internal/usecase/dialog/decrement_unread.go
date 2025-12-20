package dialog

import (
	"context"
	"fmt"

	"github.com/shaelmaar/otus-highload/social-network/counters/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/counters/internal/dto"
)

func (u *UseCases) DecrementUnreadMessages(ctx context.Context, input dto.UnreadDialogMessagesDecrement) error {
	success := true

	err := u.dialogRepo.DecrementUnreadMessages(ctx, domain.UnreadDialogMessageCountKey{
		DialogID:    input.DialogID,
		RecipientID: input.RecipientID,
	}, len(input.MessageIDs), input.IdempotencyKey)
	if err != nil {
		success = false
	}

	err = u.kafkaProducer.UnreadMessagesDecremented(ctx, dto.UnreadDialogMessagesDecremented{
		DialogID:   input.DialogID,
		MessageIDs: input.MessageIDs,
		Success:    success,
	})
	if err != nil {
		return fmt.Errorf("failed to produce unread messages decremented: %w", err)
	}

	return nil
}
