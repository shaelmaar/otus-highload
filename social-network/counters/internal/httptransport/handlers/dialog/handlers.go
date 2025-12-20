package dialog

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/counters/pkg/utils"
)

type UseCases interface {
	UnreadMessageCount(ctx context.Context, recipientID, senderID uuid.UUID) (int64, error)
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
