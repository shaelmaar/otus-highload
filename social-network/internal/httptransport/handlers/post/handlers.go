package post

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type UseCases interface {
	Create(ctx context.Context, input dto.PostCreate) (uuid.UUID, error)
	Update(ctx context.Context, input dto.PostUpdate) error
	Delete(ctx context.Context, input dto.PostDelete) error
	GetByID(ctx context.Context, id uuid.UUID) (domain.Post, error)
}

type Handlers struct {
	useCases UseCases
	logger   *zap.Logger
}

func New(useCases UseCases, logger *zap.Logger) (*Handlers, error) {
	if utils.IsNil(useCases) {
		return nil, errors.New("post use cases are nil")
	}

	if utils.IsNil(logger) {
		return nil, errors.New("logger is nil")
	}

	return &Handlers{
		useCases: useCases,
		logger:   logger,
	}, nil
}
