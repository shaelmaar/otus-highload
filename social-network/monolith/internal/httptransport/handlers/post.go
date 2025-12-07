//nolint:dupl // по какой-то причине считается дубликатом.
package handlers

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
)

// PostPostCreate создание поста (POST /post/create).
func (h *Handlers) PostPostCreate(
	ctx context.Context,
	req serverhttp.PostPostCreateRequestObject) (serverhttp.PostPostCreateResponseObject, error) {
	return h.post.Create(ctx, req)
}

// GetPostGetId получить пост по идентификатору (GET /post/get/{id}).
func (h *Handlers) GetPostGetId(
	ctx context.Context,
	req serverhttp.GetPostGetIdRequestObject) (serverhttp.GetPostGetIdResponseObject, error) {
	return h.post.GetByID(ctx, req)
}

// PutPostUpdate обновление поста (PUT /post/update).
func (h *Handlers) PutPostUpdate(
	ctx context.Context,
	req serverhttp.PutPostUpdateRequestObject) (serverhttp.PutPostUpdateResponseObject, error) {
	return h.post.Update(ctx, req)
}

// PutPostDeleteId удаление поста (PUT /post/delete/{id}).
func (h *Handlers) PutPostDeleteId(
	ctx context.Context,
	req serverhttp.PutPostDeleteIdRequestObject) (serverhttp.PutPostDeleteIdResponseObject, error) {
	return h.post.Delete(ctx, req)
}

// GetPostFeed лента постов (GET /post/feed).
func (h *Handlers) GetPostFeed(
	ctx context.Context,
	req serverhttp.GetPostFeedRequestObject) (serverhttp.GetPostFeedResponseObject, error) {
	return h.post.Feed(ctx, req)
}
