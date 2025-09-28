package friend

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

func (h *Handlers) Set(
	ctx context.Context,
	req serverhttp.PutFriendSetUserIdRequestObject) (serverhttp.PutFriendSetUserIdResponseObject, error) {
	userID, _ := auth.ExtractUserIDFromContext(ctx)

	friendID, err := uuid.Parse(req.UserId)
	if err != nil {
		//nolint:nilerr // возвращается 400 ответ.
		return serverhttp.PutFriendSetUserId400Response{}, nil
	}

	if userID == friendID {
		return serverhttp.PutFriendSetUserId400Response{}, nil
	}

	err = h.useCases.Set(ctx, dto.FriendSet{
		UserID:   userID,
		FriendID: friendID,
	})

	switch {
	case errors.Is(err, domain.ErrFriendNotFound), errors.Is(err, domain.ErrUserNotFound):
		return serverhttp.PutFriendSetUserId404Response{}, nil
	case err != nil:
		h.logger.Error("internal error", zap.Error(err))

		return serverhttp.PutFriendSetUserId500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	return serverhttp.PutFriendSetUserId200Response{}, nil
}
