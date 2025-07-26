package handlers

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type UserHandlers interface {
	Login(ctx context.Context,
		req serverhttp.PostLoginRequestObject) (serverhttp.PostLoginResponseObject, error)
	GetByID(ctx context.Context,
		req serverhttp.GetUserGetIdRequestObject) (serverhttp.GetUserGetIdResponseObject, error)
	Register(ctx context.Context,
		req serverhttp.PostUserRegisterRequestObject) (serverhttp.PostUserRegisterResponseObject, error)
	UserSearch(ctx context.Context,
		req serverhttp.GetUserSearchRequestObject) (serverhttp.GetUserSearchResponseObject, error)
}

type Handlers struct {
	user UserHandlers
}

func NewHandlers(user UserHandlers) *Handlers {
	if utils.IsNil(user) {
		panic("user handlers are nil")
	}

	return &Handlers{
		user: user,
	}
}
