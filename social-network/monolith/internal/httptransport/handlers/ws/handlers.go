package ws

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/ctxcarrier"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type Handlers struct {
	m      *melody.Melody
	logger *zap.Logger
}

func New(m *melody.Melody, logger *zap.Logger) (*Handlers, error) {
	if utils.IsNil(m) {
		return nil, errors.New("melody is nil")
	}

	if utils.IsNil(logger) {
		return nil, errors.New("logger is nil")
	}

	return &Handlers{
		m:      m,
		logger: logger,
	}, nil
}

func (h *Handlers) WS(c echo.Context) error {
	userID, _ := ctxcarrier.ExtractUserID(c.Request().Context())

	err := h.m.HandleRequestWithKeys(c.Response(), c.Request(), map[string]any{
		"user_id": userID,
	})
	if err != nil {
		h.logger.Error("failed to handle request with keys in melody", zap.Error(err))

		c.Error(err)
	}

	return nil
}
