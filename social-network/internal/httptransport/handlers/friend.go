package handlers

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
)

// PutFriendSetUserId добавление друга (PUT /friend/set/{user_id}).
func (h *Handlers) PutFriendSetUserId(
	ctx context.Context,
	req serverhttp.PutFriendSetUserIdRequestObject) (serverhttp.PutFriendSetUserIdResponseObject, error) {
	return h.friend.Set(ctx, req)
}

// PutFriendDeleteUserId удалить друга (PUT /friend/delete/{user_id}).
func (h *Handlers) PutFriendDeleteUserId(
	ctx context.Context,
	req serverhttp.PutFriendDeleteUserIdRequestObject) (serverhttp.PutFriendDeleteUserIdResponseObject, error) {
	return h.friend.Delete(ctx, req)
}
