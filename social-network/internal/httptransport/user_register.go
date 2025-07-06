package httptransport

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
)

// PostUserRegister регистрация пользователя (POST /user/register).
func (h *Handlers) PostUserRegister(
	context.Context, serverhttp.PostUserRegisterRequestObject,
) (serverhttp.PostUserRegisterResponseObject, error) {
	panic("not implemeted")
}
