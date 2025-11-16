package handlers

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
)

// PostDialogUserIdSend отправки сообшения пользователю (POST /dialog/{user_id}/send).
func (h *Handlers) PostDialogUserIdSend(
	ctx context.Context,
	req serverhttp.PostDialogUserIdSendRequestObject,
) (serverhttp.PostDialogUserIdSendResponseObject, error) {
	return h.dialog.Send(ctx, req)
}

// GetDialogUserIdList получить список сообщений диалого с пользователем (GET /dialog/{user_id}/list).
func (h *Handlers) GetDialogUserIdList(
	ctx context.Context,
	req serverhttp.GetDialogUserIdListRequestObject,
) (serverhttp.GetDialogUserIdListResponseObject, error) {
	return h.dialog.Get(ctx, req)
}
