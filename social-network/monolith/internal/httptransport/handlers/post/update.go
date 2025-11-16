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

func (h *Handlers) Update(
	ctx context.Context,
	req serverhttp.PutPostUpdateRequestObject) (serverhttp.PutPostUpdateResponseObject, error) {
	userID, _ := auth.ExtractUserIDFromContext(ctx)

	if req.Body.Id == "" || req.Body.Text == "" {
		return serverhttp.PutPostUpdate400Response{}, nil
	}

	postID, err := uuid.Parse(req.Body.Id)
	if err != nil {
		//nolint:nilerr // возвращается 400 ответ.
		return serverhttp.PutPostUpdate400Response{}, nil
	}

	err = h.useCases.Update(ctx, dto.PostUpdate{
		ID:      postID,
		Content: req.Body.Text,
		UserID:  userID,
	})

	switch {
	case errors.Is(err, domain.ErrPostNotFound):
		return serverhttp.PutPostUpdate404Response{}, nil
	case errors.Is(err, domain.ErrPostNotFoundForUser):
		return serverhttp.PutPostUpdate403Response{}, nil
	case err != nil:
		h.logger.Error("failed to update post", zap.Error(err))

		return serverhttp.PutPostUpdate500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	return serverhttp.PutPostUpdate200Response{}, nil
}
