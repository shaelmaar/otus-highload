package dialog

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/ctxcarrier"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/httptransport/handlers"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/pkg/utils"
)

func (h *Handlers) Get(
	ctx context.Context,
	req serverhttp.GetDialogUserIdListRequestObject,
) (serverhttp.GetDialogUserIdListResponseObject, error) {
	toUserID, _ := ctxcarrier.ExtractUserID(ctx)

	fromUserID, err := uuid.Parse(req.UserId)
	if err != nil {
		//nolint:nilerr // возвращается 400 ответ.
		return serverhttp.GetDialogUserIdList400Response{}, nil
	}

	if err = randomErr(); err != nil {
		h.logger.Error("internal error", zap.Error(err))

		//nolint:nilerr // возвращается 500 ответ.
		return serverhttp.GetDialogUserIdList500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	messages, err := h.useCases.GetMessagesList(ctx, dto.DialogMessagesListGet{
		From: fromUserID,
		To:   toUserID,
	})
	if err != nil {
		h.logger.Error("internal error", zap.Error(err))

		//nolint:nilerr // возвращается 500 ответ.
		return serverhttp.GetDialogUserIdList500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	dialogMessages := utils.MapSlice(messages, func(m domain.DialogMessage) serverhttp.DialogMessage {
		return serverhttp.DialogMessage{
			Id:    m.ID,
			From:  m.From.String(),
			Text:  m.Text,
			To:    m.To.String(),
			State: mapDialogMessageState(m.State),
		}
	})

	return serverhttp.GetDialogUserIdList200JSONResponse(dialogMessages), nil
}

func mapDialogMessageState(s domain.DialogMessageState) *serverhttp.DialogMessageState {
	switch s {
	case domain.DialogMessageStateSending:
		return utils.Ptr(serverhttp.Sending)
	case domain.DialogMessageStateFailed:
		return utils.Ptr(serverhttp.Failed)
	case domain.DialogMessageStateSent, domain.DialogMessageStateReading:
		return utils.Ptr(serverhttp.Sent)
	case domain.DialogMessageStateRead:
		return utils.Ptr(serverhttp.Read)
	}

	return nil
}

func randomErr() error {
	if rand.Float64() < 0.001 {
		time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond)

		return errors.New("random err")
	}

	return nil
}
