package dialog

import (
	"context"
	"fmt"

	dialogsGRPC "github.com/shaelmaar/otus-highload/social-network/gen/clientgrpc/dialogs"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
)

func (u *UseCases) CreateMessage(ctx context.Context, input dto.DialogCreateMessage) error {
	_, err := u.dialogsClient.CreateMessage(ctx, &dialogsGRPC.CreateMessageRequest{
		From: input.From.String(),
		To:   input.To.String(),
		Text: input.Text,
	})
	if err != nil {
		return fmt.Errorf("failed to create dialog message: %w", err)
	}

	return nil
}
