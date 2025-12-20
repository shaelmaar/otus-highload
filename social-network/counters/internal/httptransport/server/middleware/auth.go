package middleware

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/shaelmaar/otus-highload/social-network/counters/internal/ctxcarrier"
)

//nolint:gochecknoglobals // слайс константой не сделать.
var authSkipURLPrefixes = []string{}

const headerXUserID = "X-User-Id"

func authURLSkipper(c echo.Context) bool {
	for _, prefix := range authSkipURLPrefixes {
		if strings.HasPrefix(c.Path(), prefix) {
			return true
		}
	}

	return false
}

func auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if authURLSkipper(c) {
				return next(c)
			}

			userID := c.Request().Header.Get(headerXUserID)
			if userID == "" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error":   "unauthorized",
					"message": "x-user-id is required",
				})
			}

			userUUID, err := uuid.Parse(userID)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error":   "internal server error",
					"message": err.Error(),
				})
			}

			ctx := ctxcarrier.InjectUserID(c.Request().Context(), userUUID)

			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}
