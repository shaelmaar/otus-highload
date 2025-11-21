package dialog

import (
	"context"
	"fmt"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/dto"
)

func (u *UseCases) GetMessagesList(
	ctx context.Context, input dto.DialogMessagesListGet) ([]domain.DialogMessage, error) {
	messages, err := u.repo.GetMessagesByDialog(ctx, generateDialogID(input.From, input.To))
	if err != nil {
		return nil, fmt.Errorf("failed to get messages by dialog: %w", err)
	}

	return messages, nil
}
