package dialog

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/ctxcarrier"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/httptransport/handlers"
)

func (h *Handlers) ReadMessages(
	ctx context.Context,
	req serverhttp.PatchDialogUserIdReadMessagesMessageIdRequestObject,
) (serverhttp.PatchDialogUserIdReadMessagesMessageIdResponseObject, error) {
	toUserID, _ := ctxcarrier.ExtractUserID(ctx)

	fromUserID, err := uuid.Parse(req.UserId)
	if err != nil {
		//nolint:nilerr // возвращается 400 ответ.
		return serverhttp.PatchDialogUserIdReadMessagesMessageId400Response{}, nil
	}

	err = h.useCases.ReadMessages(ctx, dto.ReadMessages{
		From:      fromUserID,
		To:        toUserID,
		MessageID: req.MessageId,
	})
	if err != nil {
		h.logger.Error("internal error", zap.Error(err))

		return serverhttp.PatchDialogUserIdReadMessagesMessageId500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	return serverhttp.PatchDialogUserIdReadMessagesMessageId204Response{}, nil
}
