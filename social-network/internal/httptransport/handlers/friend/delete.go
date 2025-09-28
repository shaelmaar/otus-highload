package friend

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/auth"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers"
)

func (h *Handlers) Delete(ctx context.Context,
	req serverhttp.PutFriendDeleteUserIdRequestObject) (serverhttp.PutFriendDeleteUserIdResponseObject, error) {
	userID, _ := auth.ExtractUserIDFromContext(ctx)

	friendID, err := uuid.Parse(req.UserId)
	if err != nil {
		//nolint:nilerr // возвращается 400 ответ.
		return serverhttp.PutFriendDeleteUserId400Response{}, nil
	}

	err = h.useCases.Delete(ctx, dto.FriendDelete{
		UserID:   userID,
		FriendID: friendID,
	})
	if err != nil {
		h.logger.Error("internal error", zap.Error(err))

		return serverhttp.PutFriendDeleteUserId500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	return serverhttp.PutFriendDeleteUserId200Response{}, nil
}
