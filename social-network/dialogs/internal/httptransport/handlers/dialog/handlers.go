package dialog

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/pkg/utils"
)

type UseCases interface {
	CreateMessage(ctx context.Context, input dto.DialogCreateMessage) error
	GetMessagesList(
		ctx context.Context, input dto.DialogMessagesListGet) ([]domain.DialogMessage, error)
}

type Handlers struct {
	useCases UseCases
	logger   *zap.Logger
}

func New(useCases UseCases, logger *zap.Logger) (*Handlers, error) {
	if utils.IsNil(useCases) {
		return nil, errors.New("use cases are nil")
	}

	if utils.IsNil(logger) {
		return nil, errors.New("logger is nil")
	}

	return &Handlers{
		useCases: useCases,
		logger:   logger,
	}, nil
}
