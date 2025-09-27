package post

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

func (h *Handlers) GetByID(
	ctx context.Context,
	req serverhttp.GetPostGetIdRequestObject) (serverhttp.GetPostGetIdResponseObject, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		//nolint:nilerr // возвращается 400 ответ.
		return serverhttp.GetPostGetId400Response{}, nil
	}

	post, err := h.useCases.GetByID(ctx, id)

	switch {
	case errors.Is(err, domain.ErrPostNotFound):
		return serverhttp.GetPostGetId404Response{}, nil
	case err != nil:
		h.logger.Error("internal error", zap.Error(err))

		return serverhttp.GetPostGetId500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	return serverhttp.GetPostGetId200JSONResponse{
		AuthorUserId: utils.Ptr(post.AuthorUserID.String()),
		Id:           utils.Ptr(post.ID.String()),
		Text:         utils.Ptr(post.Content),
	}, nil
}
