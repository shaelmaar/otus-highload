package logger

import (
	"context"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type loggerKey struct{}

func NewEchoContext(c echo.Context, l *zap.Logger) {
	ctx := NewContext(c.Request().Context(), l)
	c.SetRequest(c.Request().WithContext(ctx))
}

func FromEchoContext(c echo.Context) *zap.Logger {
	return FromContext(c.Request().Context())
}

func NewContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, l)
}

func FromContext(ctx context.Context) *zap.Logger {
	l, ok := ctx.Value(loggerKey{}).(*zap.Logger)

	if !ok {
		return zap.L()
	}

	return l
}
