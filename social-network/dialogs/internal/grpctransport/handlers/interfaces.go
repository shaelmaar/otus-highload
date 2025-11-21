package handlers

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/dto"
)

type DialogUseCases interface {
	CreateMessage(ctx context.Context, input dto.DialogCreateMessage) error
	GetMessagesList(
		ctx context.Context, input dto.DialogMessagesListGet) ([]domain.DialogMessage, error)
}
