//nolint:dupl // по какой-то причине считается дубликатом.
package handlers

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
)

// PostLogin логин пользователя (POST /login).
func (h *Handlers) PostLogin(
	ctx context.Context, req serverhttp.PostLoginRequestObject,
) (serverhttp.PostLoginResponseObject, error) {
	return h.user.Login(ctx, req)
}

// GetUserGetId получить пользователя по id (GET /user/get/{id}).
func (h *Handlers) GetUserGetId(
	ctx context.Context, req serverhttp.GetUserGetIdRequestObject,
) (serverhttp.GetUserGetIdResponseObject, error) {
	return h.user.GetByID(ctx, req)
}

// PostUserRegister регистрация пользователя (POST /user/register).
func (h *Handlers) PostUserRegister(
	ctx context.Context, req serverhttp.PostUserRegisterRequestObject,
) (serverhttp.PostUserRegisterResponseObject, error) {
	return h.user.Register(ctx, req)
}

// GetUserSearch поиск пользователя. (GET /user/search).
func (h *Handlers) GetUserSearch(
	ctx context.Context, req serverhttp.GetUserSearchRequestObject,
) (serverhttp.GetUserSearchResponseObject, error) {
	return h.user.UserSearch(ctx, req)
}

// GetValidateToken валидация токена пользователя. (GET /validate_token).
func (h *Handlers) GetValidateToken(
	ctx context.Context, req serverhttp.GetValidateTokenRequestObject,
) (serverhttp.GetValidateTokenResponseObject, error) {
	return h.user.ValidateToken(ctx, req)
}
