package user

import (
	"context"

	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

func (h *Handlers) UserSearch(
	ctx context.Context, req serverhttp.GetUserSearchRequestObject,
) (serverhttp.GetUserSearchResponseObject, error) {
	users, err := h.useCases.Search(ctx, dto.SearchDTO{
		FirstName: req.Params.FirstName,
		LastName:  req.Params.LastName,
	})
	if err != nil {
		h.logger.Error("internal error", zap.Error(err))

		return serverhttp.GetUserSearch500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	return serverhttp.GetUserSearch200JSONResponse(utils.MapSlice(users, parseUser)), nil
}
