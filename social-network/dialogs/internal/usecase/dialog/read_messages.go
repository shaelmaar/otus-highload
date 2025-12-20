package dialog

import (
	"context"
	"fmt"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/dto"
)

func (u *UseCases) ReadMessages(ctx context.Context, input dto.ReadMessages) error {
	dialogID := generateDialogID(input.From, input.To)

	ids, err := u.repo.MarkMessagesAsReading(ctx, dialogID, input.To, input.MessageID)
	if err != nil {
		return fmt.Errorf("failed to mark messages as reading: %w", err)
	}

	err = u.kafkaProducer.MessagesRead(ctx, dto.MessagesReadEvent{
		MessageIDs: ids,
		DialogID:   dialogID,
		To:         input.To,
	})
	if err != nil {
		return fmt.Errorf("failed to produce messages read event: %w", err)
	}

	return nil
}
