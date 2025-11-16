package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/server/logger"
)

type zapLoggerConfig struct {
	ServiceName string
	Skipper     middleware.Skipper
}

//nolint:gochecknoglobals
var headerLogFields = []string{
	"x-device-id", "x-device-platform", "x-app-version", "x-username", "x-keycloak-id", "x-device-tag", "x-request-id",
	"user-agent", "x-user-agent", "x-client-name",
}

func zapLoggerMiddleware(config zapLoggerConfig, l *zap.Logger) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			startTime := time.Now()

			fields := make([]zap.Field, 0, 6+len(headerLogFields)) //nolint:mnd

			fields = append(
				fields,
				zap.String("method", c.Request().Method),
				zap.String("request_uri", c.Request().URL.RequestURI()),
				zap.String("ip", c.RealIP()),
				zap.String("component", config.ServiceName),
			)

			for _, fieldName := range headerLogFields {
				if fieldValue := c.Request().Header.Get(fieldName); fieldValue != "" {
					fields = append(fields, zap.String(fieldName, fieldValue))
				}
			}

			logger.NewEchoContext(c, l.With(fields...))

			err := next(c)

			statusCode := c.Response().Status

			if errors.Is(c.Request().Context().Err(), context.Canceled) {
				statusCode = 499
			} else if err != nil {
				statusCode = getEchoErrorStatusCode(err)
			}

			errField := zap.Skip()
			if err != nil {
				errField = zap.Error(err)
			}

			logger.FromEchoContext(c).Info(
				"request",
				errField,
				zap.Int("status_code", statusCode),
				zap.Duration("duration", time.Since(startTime)),
			)

			return err
		}
	}
}

func getEchoErrorStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	var httpErr *echo.HTTPError

	if errors.As(err, &httpErr) {
		if httpErr.Internal != nil {
			var innerErr *echo.HTTPError

			if errors.As(httpErr.Internal, &innerErr) {
				return innerErr.Code
			}
		}

		return httpErr.Code
	}

	return http.StatusInternalServerError
}
