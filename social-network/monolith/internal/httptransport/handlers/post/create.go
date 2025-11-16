package post

import (
	"context"

	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/auth"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers"
)

func (h *Handlers) Create(
	ctx context.Context,
	req serverhttp.PostPostCreateRequestObject) (serverhttp.PostPostCreateResponseObject, error) {
	userID, _ := auth.ExtractUserIDFromContext(ctx)

	if req.Body.Text == "" {
		return serverhttp.PostPostCreate400Response{}, nil
	}

	id, err := h.useCases.Create(ctx, dto.PostCreate{
		Content: req.Body.Text,
		UserID:  userID,
	})
	if err != nil {
		h.logger.Error("internal error", zap.Error(err))

		return serverhttp.PostPostCreate500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	return serverhttp.PostPostCreate200JSONResponse(id.String()), nil
}
