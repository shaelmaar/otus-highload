package handlers

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/gen/servergrpc"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/dto"
	grpcServer "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/grpctransport/server"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/pkg/utils"
)

func (h *Handlers) GetDialogMessages(
	ctx context.Context, req *servergrpc.GetDialogMessagesRequest) (*servergrpc.GetDialogMessagesReply, error) {
	from, err := uuid.Parse(req.From)
	if err != nil {
		return nil, grpcServer.GRPCValidationError(servergrpc.GetDialogMessagesReply_VALIDATION_ERROR, err)
	}

	to, err := uuid.Parse(req.To)
	if err != nil {
		return nil, grpcServer.GRPCValidationError(servergrpc.GetDialogMessagesReply_VALIDATION_ERROR, err)
	}

	messages, err := h.dialogUseCases.GetMessagesList(ctx, dto.DialogMessagesListGet{
		From: from,
		To:   to,
	})
	if err != nil {
		return nil, grpcServer.GRPCUnknownError(servergrpc.GetDialogMessagesReply_UNSPECIFIED, err)
	}

	return &servergrpc.GetDialogMessagesReply{
		Messages: utils.MapSlice(messages, func(m domain.DialogMessage) *servergrpc.DialogMessage {
			return &servergrpc.DialogMessage{
				From:      m.From.String(),
				To:        m.To.String(),
				Text:      m.Text,
				CreatedAt: timestamppb.New(m.CreatedAt),
			}
		}),
	}, nil
}
