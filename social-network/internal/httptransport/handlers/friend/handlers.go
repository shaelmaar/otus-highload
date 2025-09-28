package friend

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type UseCases interface {
	Set(ctx context.Context, input dto.FriendSet) error
	Delete(ctx context.Context, input dto.FriendDelete) error
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
