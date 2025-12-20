package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/shaelmaar/otus-highload/social-network/counters/internal/ctxcarrier"
)

const (
	headerXRequestID = "X-Request-Id"
)

type requestIDConfig struct {
	Skipper middleware.Skipper
}

func requestIDWithConfig(config requestIDConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			rid := c.Request().Header.Get(headerXRequestID)
			if rid == "" {
				rid = uuid.New().String()

				c.Request().Header.Set(headerXRequestID, rid)
			}

			c.Response().Header().Set(headerXRequestID, rid)

			ctx := ctxcarrier.InjectRequestID(c.Request().Context(), rid)

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
