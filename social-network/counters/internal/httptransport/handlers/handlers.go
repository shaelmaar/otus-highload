package handlers

import (
	"context"
	"errors"

	"github.com/shaelmaar/otus-highload/social-network/counters/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/counters/pkg/utils"
)

type DialogHandlers interface {
	GetUnread(
		ctx context.Context, req serverhttp.GetDialogUserIdUnreadRequestObject,
	) (serverhttp.GetDialogUserIdUnreadResponseObject, error)
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
