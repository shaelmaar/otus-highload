package loadtest

import (
	"context"
	"errors"

	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type UseCases interface {
	Write(ctx context.Context, value string) error
}

type Handlers struct {
	useCases UseCases
	logger   *zap.Logger
}

func New(useCases UseCases, logger *zap.Logger) (*Handlers, error) {
	if utils.IsNil(useCases) {
		return nil, errors.New("use cases is nil")
	}

	if utils.IsNil(logger) {
		return nil, errors.New("logger is nil")
	}

	return &Handlers{
		useCases: useCases,
		logger:   logger,
	}, nil
}

func (h *Handlers) Write(
	ctx context.Context,
	req serverhttp.PostLoadtestWriteRequestObject,
) (serverhttp.PostLoadtestWriteResponseObject, error) {
	err := h.useCases.Write(ctx, req.Body.Value)
	if err != nil {
		h.logger.Error("failed to write load test", zap.Error(err))

		return serverhttp.PostLoadtestWrite500JSONResponse{
			N5xxJSONResponse: handlers.Simple500JSONResponse(""),
		}, nil
	}

	return serverhttp.PostLoadtestWrite204Response{}, nil
}
