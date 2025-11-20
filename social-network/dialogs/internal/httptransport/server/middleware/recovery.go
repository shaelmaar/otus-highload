package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/httptransport/server/logger"
)

func recovery() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}

					logger.FromEchoContext(c).Error(
						"[PANIC RECOVER]",
						zap.Error(err),
						zap.Stack("stack"),
					)

					c.Error(err)
				}
			}()

			return next(c)
		}
	}
}
