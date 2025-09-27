package post

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/auth"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers"
)

func (h *Handlers) Delete(
	ctx context.Context,
	req serverhttp.PutPostDeleteIdRequestObject) (serverhttp.PutPostDeleteIdResponseObject, error) {
	userID, _ := auth.ExtractUserIDFromContext(ctx)

	if req.Id == "" {
		return serverhttp.PutPostDeleteId400Response{}, nil
	}

	postID, err := uuid.Parse(req.Id)
	if err != nil {
		//nolint:nilerr // возвращается 400 ответ.
		return serverhttp.PutPostDeleteId400Response{}, nil
	}

	err = h.useCases.Delete(ctx, dto.PostDelete{
		ID:     postID,
		UserID: userID,
	})

	switch {
	case errors.Is(err, domain.ErrPostNotFound):
		return serverhttp.PutPostDeleteId404Response{}, nil
	case errors.Is(err, domain.ErrPostNotFoundForUser):
		return serverhttp.PutPostDeleteId403Response{}, nil
	case err != nil:
		h.logger.Error("failed to delete post", zap.Error(err))

		return serverhttp.PutPostDeleteId500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	return serverhttp.PutPostDeleteId200Response{}, nil
}
