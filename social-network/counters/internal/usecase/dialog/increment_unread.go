package dialog

import (
	"context"
	"fmt"

	"github.com/shaelmaar/otus-highload/social-network/counters/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/counters/internal/dto"
)

func (u *UseCases) IncrementUnreadMessages(ctx context.Context, input dto.UnreadDialogMessagesIncrement) error {
	success := true

	err := u.dialogRepo.IncrementUnreadMessages(ctx, domain.UnreadDialogMessageCountKey{
		DialogID:    input.DialogID,
		RecipientID: input.RecipientID,
	}, input.IdempotencyKey)
	if err != nil {
		success = false
	}

	err = u.kafkaProducer.UnreadMessagesIncremented(ctx, dto.UnreadDialogMessagesIncremented{
		DialogID:  input.DialogID,
		MessageID: input.MessageID,
		Success:   success,
	})
	if err != nil {
		return fmt.Errorf("failed to produce unread messages incremented: %w", err)
	}

	return nil
}
