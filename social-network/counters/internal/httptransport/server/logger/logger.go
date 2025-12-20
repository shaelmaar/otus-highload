package logger

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/counters/internal/ctxcarrier"
)

type HTTPServerLogger struct {
	logger *zap.Logger
}

func (l *HTTPServerLogger) Write(p []byte) (int, error) {
	l.logger.Error(string(p))

	return len(p), nil
}

func NewHTTPServerLogger(l *zap.Logger) *HTTPServerLogger {
	return &HTTPServerLogger{l}
}

func NewEchoContext(c echo.Context, l *zap.Logger) {
	ctx := ctxcarrier.InjectLogger(c.Request().Context(), l)
	c.SetRequest(c.Request().WithContext(ctx))
}

func FromEchoContext(c echo.Context) *zap.Logger {
	return ctxcarrier.ExtractLogger(c.Request().Context())
}
