package dialog

import (
	"context"
	"fmt"

	dialogsGRPC "github.com/shaelmaar/otus-highload/social-network/gen/clientgrpc/dialogs"
	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

func (u *UseCases) GetMessagesList(
	ctx context.Context, input dto.DialogMessagesListGet) ([]domain.DialogMessage, error) {
	reply, err := u.dialogsClient.GetDialogMessages(ctx, &dialogsGRPC.GetDialogMessagesRequest{
		From: input.From.String(),
		To:   input.To.String(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get dialog messages: %w", err)
	}

	return utils.MapSlice(reply.Messages, func(m *dialogsGRPC.DialogMessage) domain.DialogMessage {
		return domain.DialogMessage{
			From:      m.From,
			To:        m.To,
			Text:      m.Text,
			CreatedAt: m.CreatedAt.AsTime(),
		}
	}), nil
}
