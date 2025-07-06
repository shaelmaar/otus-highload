package httptransport

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
)

// PostLogin логин пользователя (POST /login).
func (h *Handlers) PostLogin(
	context.Context, serverhttp.PostLoginRequestObject,
) (serverhttp.PostLoginResponseObject, error) {
	panic("not implemeted")
}
