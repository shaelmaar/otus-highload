package user

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
	Login(ctx context.Context, dto dto.LoginDTO) (domain.UserToken, error)
	Register(ctx context.Context, dto dto.RegisterDTO) (domain.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (domain.User, error)
}

type Handlers struct {
	useCases UseCases
	logger   *zap.Logger
}

func NewHandlers(useCases UseCases, logger *zap.Logger) (*Handlers, error) {
	if utils.IsNil(useCases) {
		return nil, errors.New("user use cases are nil")
	}

	if utils.IsNil(logger) {
		return nil, errors.New("logger is nil")
	}

	return &Handlers{
		useCases: useCases,
		logger:   logger,
	}, nil
}
