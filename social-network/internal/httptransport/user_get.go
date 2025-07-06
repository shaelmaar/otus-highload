package httptransport

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
)

// GetUserGetId получить пользователя по id (GET /user/get/{id}).
func (h *Handlers) GetUserGetId(
	context.Context, serverhttp.GetUserGetIdRequestObject,
) (serverhttp.GetUserGetIdResponseObject, error) {
	panic("not implemeted")
}
