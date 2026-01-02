package middleware

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

var httpServerRequest *prometheus.SummaryVec

type metricsConfig struct {
	Skipper     middleware.Skipper
	ServiceName string
}

func metricsMiddleware(config metricsConfig) echo.MiddlewareFunc {
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultSkipper
	}

	httpServerRequest = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace:  "http_server",
			Name:       "request",
			Help:       "The duration of server requests in seconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.75: 0.025, 0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
			ConstLabels: map[string]string{
				"service_name": config.ServiceName,
			},
		},
		[]string{
			"http_method",
			"http_path",
			"http_status_code",
		},
	)

	prometheus.MustRegister(httpServerRequest)

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			startTime := time.Now()

			path := c.Path()

			err := next(c)

			statusCode := c.Response().Status

			if errors.Is(c.Request().Context().Err(), context.Canceled) {
				statusCode = 499
			} else if err != nil {
				statusCode = getEchoErrorStatusCode(err)
			}

			httpServerRequest.With(prometheus.Labels{
				"http_method":      strings.ToLower(c.Request().Method),
				"http_path":        path,
				"http_status_code": strconv.Itoa(statusCode),
			}).Observe(time.Since(startTime).Seconds())

			return err
		}
	}
}
