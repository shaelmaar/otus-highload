package dialog

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/counters/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/counters/internal/ctxcarrier"
	"github.com/shaelmaar/otus-highload/social-network/counters/internal/httptransport/handlers"
)

func (h *Handlers) GetUnread(
	ctx context.Context, req serverhttp.GetDialogUserIdUnreadRequestObject,
) (serverhttp.GetDialogUserIdUnreadResponseObject, error) {
	recipientID, _ := ctxcarrier.ExtractUserID(ctx)

	fromUserID, err := uuid.Parse(req.UserId)
	if err != nil {
		//nolint:nilerr // возвращается 400 ответ.
		return serverhttp.GetDialogUserIdUnread400Response{}, nil
	}

	res, err := h.useCases.UnreadMessageCount(ctx, recipientID, fromUserID)
	if err != nil {
		h.logger.Error("internal error", zap.Error(err))
		//nolint:nilerr // возвращается 500 ответ.
		return serverhttp.GetDialogUserIdUnread500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	return serverhttp.GetDialogUserIdUnread200JSONResponse(res), nil
}
