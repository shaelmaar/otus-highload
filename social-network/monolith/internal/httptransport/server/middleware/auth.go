package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/shaelmaar/otus-highload/social-network/internal/ctxcarrier"
)

//nolint:gochecknoglobals // слайс константой не сделать.
var authSkipURLPrefixes = []string{
	"/login",
	"/user/register",
	"/user/get",
	"/user/search",
	"/post/get",
}

type authConfig struct {
	tokenValidator func(string) (string, error)
}

const authHeader = "Authorization"

func authURLSkipper(c echo.Context) bool {
	fmt.Println(c.Path())

	for _, prefix := range authSkipURLPrefixes {
		if strings.HasPrefix(c.Path(), prefix) {
			return true
		}
	}

	return false
}

func authWithConfig(cfg authConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if authURLSkipper(c) {
				return next(c)
			}

			token, err := extractTokenFromHeader(c)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error":   "unauthorized",
					"message": err.Error(),
				})
			}

			userID, err := cfg.tokenValidator(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error":   "unauthorized",
					"message": err.Error(),
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

func extractTokenFromHeader(c echo.Context) (string, error) {
	header := c.Request().Header.Get(authHeader)
	if header == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
	}

	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return "", echo.NewHTTPError(http.StatusUnauthorized,
			"Invalid authorization header format")
	}

	if parts[0] != "Bearer" {
		return "", echo.NewHTTPError(http.StatusUnauthorized,
			"Authorization header must start with 'Bearer'")
	}

	return parts[1], nil
}
