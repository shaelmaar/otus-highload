package handlers

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/counters/gen/serverhttp"
)

// GetDialogUserIdUnread кол-во непрочитанных сообщений в диалоге.
func (h *Handlers) GetDialogUserIdUnread(
	ctx context.Context, req serverhttp.GetDialogUserIdUnreadRequestObject,
) (serverhttp.GetDialogUserIdUnreadResponseObject, error) {
	return h.dialog.GetUnread(ctx, req)
}
