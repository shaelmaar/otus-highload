package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers"
)

func (h *Handlers) Login(
	ctx context.Context, req serverhttp.PostLoginRequestObject,
) (serverhttp.PostLoginResponseObject, error) {
	if req.Body.Id == nil || req.Body.Password == nil {
		return serverhttp.PostLogin400Response{}, nil
	}

	userID, err := uuid.Parse(*req.Body.Id)
	if err != nil {
		//nolint:nilerr // возвращается 400 ответ.
		return serverhttp.PostLogin400Response{}, nil
	}

	token, err := h.useCases.Login(ctx, dto.Login{
		UserID:   userID,
		Password: *req.Body.Password,
	})

	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return serverhttp.PostLogin404Response{}, nil
	case errors.Is(err, domain.ErrInvalidCredentials):
		return serverhttp.PostLogin400Response{}, nil
	case err != nil:
		h.logger.Error("internal error", zap.Error(err))

		return serverhttp.PostLogin500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	return serverhttp.PostLogin200JSONResponse{Token: &token}, nil
}
