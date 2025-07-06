package httptransport

import (
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type Handlers struct {
	logger *zap.Logger
}

func NewHandlers(logger *zap.Logger) *Handlers {
	if utils.IsNil(logger) {
		panic("logger is nil")
	}

	return &Handlers{
		logger: logger,
	}
}
