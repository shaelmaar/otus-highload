package user

import (
	"context"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/internal/ctxcarrier"
)

func (h *Handlers) ValidateToken(
	ctx context.Context, _ serverhttp.GetValidateTokenRequestObject,
) (serverhttp.GetValidateTokenResponseObject, error) {
	userID, _ := ctxcarrier.ExtractUserID(ctx)

	return serverhttp.GetValidateToken204Response{
		Headers: serverhttp.GetValidateToken204ResponseHeaders{
			XUserId: userID.String(),
		},
	}, nil
}
