package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

func (h *Handlers) GetByID(
	ctx context.Context,
	req serverhttp.GetUserGetIdRequestObject,
) (serverhttp.GetUserGetIdResponseObject, error) {
	userID, err := uuid.Parse(req.Id)
	if err != nil {
		//nolint:nilerr // пустой ответ в контрактах.
		return serverhttp.GetUserGetId400Response{}, nil
	}

	user, err := h.useCases.GetByID(ctx, userID)

	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return serverhttp.GetUserGetId404Response{}, nil
	case err != nil:
		h.logger.Error("internal error", zap.Error(err))

		return serverhttp.GetUserGetId500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	return serverhttp.GetUserGetId200JSONResponse(parseUser(user)), nil
}

func parseUser(user domain.User) serverhttp.User {
	return serverhttp.User{
		Biography: utils.Ptr(user.Biography),
		Birthdate: utils.Ptr(openapitypes.Date{
			Time: user.BirthDate,
		}),
		City:       utils.Ptr(user.City),
		FirstName:  utils.Ptr(user.FirstName),
		Id:         utils.Ptr(user.ID.String()),
		SecondName: utils.Ptr(user.SecondName),
	}
}
