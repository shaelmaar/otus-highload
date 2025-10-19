package dialog

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/auth"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers"
)

func (h *Handlers) Send(
	ctx context.Context,
	req serverhttp.PostDialogUserIdSendRequestObject,
) (serverhttp.PostDialogUserIdSendResponseObject, error) {
	fromUserID, _ := auth.ExtractUserIDFromContext(ctx)

	toUserID, err := uuid.Parse(req.UserId)
	if err != nil {
		//nolint:nilerr // возвращается 400 ответ.
		return serverhttp.PostDialogUserIdSend400Response{}, nil
	}

	err = h.useCases.CreateMessage(ctx, dto.DialogCreateMessage{
		From: fromUserID,
		To:   toUserID,
		Text: req.Body.Text,
		Time: time.Now(),
	})
	if err != nil {
		//nolint:nilerr // возвращается 500 ответ.
		return serverhttp.PostDialogUserIdSend500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	return serverhttp.PostDialogUserIdSend200Response{}, nil
}
