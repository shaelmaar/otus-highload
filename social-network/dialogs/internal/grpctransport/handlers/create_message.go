package handlers

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/gen/servergrpc"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/dto"
	grpcServer "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/grpctransport/server"
)

func (h *Handlers) CreateMessage(
	ctx context.Context, req *servergrpc.CreateMessageRequest) (*servergrpc.CreateMessageReply, error) {
	if req.Text == "" {
		return nil, grpcServer.GRPCValidationError(servergrpc.CreateMessageReply_VALIDATION_ERROR,
			errors.New("text is required"))
	}

	from, err := uuid.Parse(req.From)
	if err != nil {
		return nil, grpcServer.GRPCValidationError(servergrpc.CreateMessageReply_VALIDATION_ERROR, err)
	}

	to, err := uuid.Parse(req.To)
	if err != nil {
		return nil, grpcServer.GRPCValidationError(servergrpc.CreateMessageReply_VALIDATION_ERROR, err)
	}

	err = h.dialogUseCases.CreateMessage(ctx, dto.DialogCreateMessage{
		From: from,
		To:   to,
		Text: req.Text,
		Time: time.Now(),
	})
	if err != nil {
		return nil, grpcServer.GRPCUnknownError(servergrpc.CreateMessageReply_UNSPECIFIED, err)
	}

	return &servergrpc.CreateMessageReply{}, nil
}
