package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
)

func (h *Handlers) Login(
	ctx context.Context, req serverhttp.PostLoginRequestObject,
) (serverhttp.PostLoginResponseObject, error) {
	if req.Body.Id == nil || req.Body.Password == nil {
		return serverhttp.PostLogin400Response{}, nil
	}

	userID, err := uuid.Parse(*req.Body.Id)
	if err != nil {
		//nolint:nilerr // пустой ответ в контрактах.
		return serverhttp.PostLogin400Response{}, nil
	}

	token, err := h.useCases.Login(ctx, dto.LoginDTO{
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
			N5xxJSONResponse: serverhttp.N5xxJSONResponse{
				Body: struct {
					Code      *int    `json:"code,omitempty"`
					Message   string  `json:"message"`
					RequestId *string `json:"request_id,omitempty"`
				}{
					Code:      nil,
					Message:   "Внутренняя ошибка сервера",
					RequestId: nil,
				},
				Headers: serverhttp.N5xxResponseHeaders{
					RetryAfter: 0,
				},
			},
		}, nil
	}

	return serverhttp.PostLogin200JSONResponse{Token: &token.Token}, nil
}
