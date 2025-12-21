package handlers

import (
	"context"
	"errors"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/pkg/utils"
)

type DialogHandlers interface {
	Send(ctx context.Context,
		req serverhttp.PostDialogUserIdSendRequestObject) (serverhttp.PostDialogUserIdSendResponseObject, error)
	Get(ctx context.Context,
		req serverhttp.GetDialogUserIdListRequestObject) (serverhttp.GetDialogUserIdListResponseObject, error)
	ReadMessages(
		ctx context.Context,
		req serverhttp.PatchDialogUserIdReadMessagesMessageIdRequestObject,
	) (serverhttp.PatchDialogUserIdReadMessagesMessageIdResponseObject, error)
}

type Handlers struct {
	dialog DialogHandlers
}

func New(
	dialog DialogHandlers,
) (*Handlers, error) {
	if utils.IsNil(dialog) {
		return nil, errors.New("dialog handlers are nil")
	}

	return &Handlers{
		dialog: dialog,
	}, nil
}
