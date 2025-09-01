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

type LoadTestHandlers interface {
	Write(ctx context.Context,
		req serverhttp.PostLoadtestWriteRequestObject,
	) (serverhttp.PostLoadtestWriteResponseObject, error)
}

type Handlers struct {
	user     UserHandlers
	loadTest LoadTestHandlers
}

func NewHandlers(user UserHandlers, loadTest LoadTestHandlers) *Handlers {
	if utils.IsNil(user) {
		panic("user handlers are nil")
	}

	if utils.IsNil(loadTest) {
		panic("load test handlers are nil")
	}

	return &Handlers{
		user:     user,
		loadTest: loadTest,
	}
}
