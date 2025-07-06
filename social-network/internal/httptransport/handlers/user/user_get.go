package user

import (
	"context"
	"errors"

	"github.com/google/uuid"
	openapitypes "github.com/oapi-codegen/runtime/types"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

func (h *Handlers) GetByID(
	ctx context.Context,
	req serverhttp.GetUserGetIdRequestObject,
) (serverhttp.GetUserGetIdResponseObject, error) {
	userID, err := uuid.Parse(req.Id)
	if err != nil {
		return serverhttp.GetUserGetId400Response{}, nil
	}

	user, err := h.useCases.GetByID(ctx, userID)

	switch {
	case errors.Is(err, domain.ErrUserNotFound):
		return serverhttp.GetUserGetId404Response{}, nil
	case err != nil:
		h.logger.Error("internal error", zap.Error(err))

		return serverhttp.GetUserGetId500JSONResponse{
			N5xxJSONResponse: serverhttp.N5xxJSONResponse{
				Body: struct {
					Code      *int    `json:"code,omitempty"`
					Message   string  `json:"message"`
					RequestId *string `json:"request_id,omitempty"`
				}{
					Message: "Внутренняя ошибка сервера",
				},
				Headers: serverhttp.N5xxResponseHeaders{},
			},
		}, nil
	}

	return serverhttp.GetUserGetId200JSONResponse{
		Biography: utils.Ptr(user.Biography),
		Birthdate: utils.Ptr(openapitypes.Date{
			Time: user.BirthDate,
		}),
		City:       utils.Ptr(user.City),
		FirstName:  utils.Ptr(user.FirstName),
		Id:         utils.Ptr(user.ID.String()),
		SecondName: utils.Ptr(user.SecondName),
	}, nil
}
