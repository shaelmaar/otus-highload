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

type PostHandlers interface {
	Create(ctx context.Context,
		req serverhttp.PostPostCreateRequestObject) (serverhttp.PostPostCreateResponseObject, error)
	GetByID(ctx context.Context,
		req serverhttp.GetPostGetIdRequestObject) (serverhttp.GetPostGetIdResponseObject, error)
	Update(ctx context.Context,
		req serverhttp.PutPostUpdateRequestObject) (serverhttp.PutPostUpdateResponseObject, error)
	Delete(ctx context.Context,
		req serverhttp.PutPostDeleteIdRequestObject) (serverhttp.PutPostDeleteIdResponseObject, error)
	Feed(ctx context.Context,
		req serverhttp.GetPostFeedRequestObject) (serverhttp.GetPostFeedResponseObject, error)
}

type FriendHandlers interface {
	Set(ctx context.Context,
		req serverhttp.PutFriendSetUserIdRequestObject) (serverhttp.PutFriendSetUserIdResponseObject, error)
	Delete(ctx context.Context,
		req serverhttp.PutFriendDeleteUserIdRequestObject) (serverhttp.PutFriendDeleteUserIdResponseObject, error)
}

type LoadTestHandlers interface {
	Write(ctx context.Context,
		req serverhttp.PostLoadtestWriteRequestObject,
	) (serverhttp.PostLoadtestWriteResponseObject, error)
}

type Handlers struct {
	user     UserHandlers
	post     PostHandlers
	friend   FriendHandlers
	loadTest LoadTestHandlers
}

func NewHandlers(
	user UserHandlers,
	post PostHandlers,
	friend FriendHandlers,
	loadTest LoadTestHandlers,
) *Handlers {
	if utils.IsNil(user) {
		panic("user handlers are nil")
	}

	if utils.IsNil(post) {
		panic("post handlers are nil")
	}

	if utils.IsNil(friend) {
		panic("friend handlers are nil")
	}

	if utils.IsNil(loadTest) {
		panic("load test handlers are nil")
	}

	return &Handlers{
		user:     user,
		post:     post,
		friend:   friend,
		loadTest: loadTest,
	}
}
