package handlers

import (
	"errors"

	"github.com/labstack/echo/v4"

	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

type WSHandler interface {
	WS(c echo.Context) error
}

type WSHandlers struct {
	ws WSHandler
}

func NewWSHandlers(ws WSHandler) (*WSHandlers, error) {
	if utils.IsNil(ws) {
		return nil, errors.New("ws handler is nil")
	}

	return &WSHandlers{ws: ws}, nil
}

func (h *WSHandlers) HandleWebSocket(c echo.Context) error {
	return h.ws.WS(c)
}
