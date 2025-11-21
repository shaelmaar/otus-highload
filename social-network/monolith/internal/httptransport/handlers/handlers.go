package handlers

import (
	"context"
	"errors"

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

type DialogHandlers interface {
	Send(ctx context.Context,
		req serverhttp.PostDialogUserIdSendRequestObject) (serverhttp.PostDialogUserIdSendResponseObject, error)
	Get(ctx context.Context,
		req serverhttp.GetDialogUserIdListRequestObject) (serverhttp.GetDialogUserIdListResponseObject, error)
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
	dialog   DialogHandlers
	loadTest LoadTestHandlers
}

func NewHandlers(
	user UserHandlers,
	post PostHandlers,
	friend FriendHandlers,
	dialog DialogHandlers,
	loadTest LoadTestHandlers,
) (*Handlers, error) {
	if utils.IsNil(user) {
		return nil, errors.New("user handlers are nil")
	}

	if utils.IsNil(post) {
		return nil, errors.New("post handlers are nil")
	}

	if utils.IsNil(friend) {
		return nil, errors.New("friend handlers are nil")
	}

	if utils.IsNil(dialog) {
		return nil, errors.New("dialog handlers are nil")
	}

	if utils.IsNil(loadTest) {
		return nil, errors.New("load test handlers are nil")
	}

	return &Handlers{
		user:     user,
		post:     post,
		friend:   friend,
		dialog:   dialog,
		loadTest: loadTest,
	}, nil
}
