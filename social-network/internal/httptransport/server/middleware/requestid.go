package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const headerXRequestID = "X-Request-Id"

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

			return next(c)
		}
	}
}
