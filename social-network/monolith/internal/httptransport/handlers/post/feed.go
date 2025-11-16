package post

import (
	"context"

	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/auth"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

func (h *Handlers) Feed(
	ctx context.Context,
	req serverhttp.GetPostFeedRequestObject) (serverhttp.GetPostFeedResponseObject, error) {
	userID, _ := auth.ExtractUserIDFromContext(ctx)

	const (
		maxOffset = 1000
		maxLimit  = 1000
	)

	offset := 0
	limit := 10

	if req.Params.Offset != nil {
		offset = int(*req.Params.Offset)
	}

	if req.Params.Limit != nil {
		limit = int(*req.Params.Limit)
	}

	switch {
	case offset >= maxOffset:
		return serverhttp.GetPostFeed200JSONResponse{}, nil
	case offset+limit >= maxLimit:
		limit = maxLimit - offset
	case limit > maxLimit:
		limit = maxLimit
	}

	posts, err := h.useCases.GetPostFeed(ctx, dto.GetPostFeed{
		UserID: userID,
		Offset: offset,
		Limit:  limit,
	})
	if err != nil {
		h.logger.Error("internal error", zap.Error(err))

		return serverhttp.GetPostFeed500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, err
	}

	return serverhttp.GetPostFeed200JSONResponse(utils.MapSlice(posts, func(post domain.Post) serverhttp.Post {
		return serverhttp.Post{
			AuthorUserId: utils.Ptr(post.AuthorUserID.String()),
			Id:           utils.Ptr(post.ID.String()),
			Text:         utils.Ptr(post.Content),
		}
	})), nil
}
