package dialog

import (
	"context"

	"github.com/google/uuid"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/internal/ctxcarrier"
	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

func (h *Handlers) Get(
	ctx context.Context,
	req serverhttp.GetDialogUserIdListRequestObject,
) (serverhttp.GetDialogUserIdListResponseObject, error) {
	fromUserID, _ := ctxcarrier.ExtractUserID(ctx)

	toUserID, err := uuid.Parse(req.UserId)
	if err != nil {
		//nolint:nilerr // возвращается 400 ответ.
		return serverhttp.GetDialogUserIdList400Response{}, nil
	}

	messages, err := h.useCases.GetMessagesList(ctx, dto.DialogMessagesListGet{
		From: fromUserID,
		To:   toUserID,
	})
	if err != nil {
		//nolint:nilerr // возвращается 500 ответ.
		return serverhttp.GetDialogUserIdList500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	dialogMessages := utils.MapSlice(messages, func(m domain.DialogMessage) serverhttp.DialogMessage {
		return serverhttp.DialogMessage{
			From: m.From,
			Text: m.Text,
			To:   m.To,
		}
	})

	return serverhttp.GetDialogUserIdList200JSONResponse(dialogMessages), nil
}
