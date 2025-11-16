package middleware

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type Options struct {
	ServiceName    string
	Logger         *zap.Logger
	TokenValidator func(tokenString string) (string, error)

	RequestIDGenerator func() string

	RequestIDSkipper func(echo.Context) bool
	MetricsSkipper   func(echo.Context) bool
	TraceSkipper     func(echo.Context) bool
	LoggerSkipper    func(echo.Context) bool
}

func Use(e *echo.Echo, opt *Options) {
	e.Use(
		recovery(),
		requestIDWithConfig(requestIDConfig{Skipper: opt.RequestIDSkipper}),
		authWithConfig(authConfig{tokenValidator: opt.TokenValidator}),
		zapLoggerMiddleware(zapLoggerConfig{
			ServiceName: opt.ServiceName,
			Skipper:     opt.LoggerSkipper,
		}, opt.Logger),
	)
}
