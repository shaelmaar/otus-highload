package debugserver

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func New() *echo.Echo {
	debugSrv := echo.New()

	debugSrv.GET("/metrics", echo.WrapHandler(promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer, promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{EnableOpenMetrics: false}, //nolint:exhaustruct // остальное по умолчанию.
		),
	)))

	debugSrv.GET("/ready", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	return debugSrv
}
