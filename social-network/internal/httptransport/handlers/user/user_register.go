package user

import (
	"context"

	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

func (h *Handlers) Register(
	ctx context.Context, req serverhttp.PostUserRegisterRequestObject,
) (serverhttp.PostUserRegisterResponseObject, error) {
	if req.Body.Password == nil || req.Body.FirstName == nil || req.Body.SecondName == nil ||
		req.Body.Birthdate == nil {
		return serverhttp.PostUserRegister400Response{}, nil
	}

	user, err := h.useCases.Register(ctx, dto.RegisterDTO{
		Password:   *req.Body.Password,
		Name:       *req.Body.FirstName,
		SecondName: *req.Body.SecondName,
		BirthDate:  req.Body.Birthdate.Time,
		Biography:  utils.UnPtr(req.Body.Biography),
		City:       utils.UnPtr(req.Body.City),
	})
	if err != nil {
		h.logger.Error("internal error", zap.Error(err))

		return serverhttp.PostUserRegister500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	return serverhttp.PostUserRegister200JSONResponse{UserId: utils.Ptr(user.ID.String())}, nil
}
