package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/httptransport/server/logger"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/httptransport/server/middleware"
)

type Server struct {
	Echo   *echo.Echo
	Logger *zap.Logger
}

type AuthService interface {
	ValidateToken(ctx context.Context, tokenString string) (string, error)
}

type Options struct {
	ServiceName string
	Debug       bool
	AuthService AuthService
	Logger      *zap.Logger
}

type optionalOptions struct {
	RequestIDSkipper func(echo.Context) bool
	MetricsSkipper   func(echo.Context) bool
	TraceSkipper     func(echo.Context) bool
	LoggerSkipper    func(echo.Context) bool

	Middleware []echo.MiddlewareFunc
}

//nolint:gochecknoglobals
var skipURLPrefixList = []string{
	"/metrics",
	"/readiness",
	"/liveness",
	"/debug",
}

type OptionalFunc func(opt *optionalOptions)

func WithRequestIDSkipper(fn func(echo.Context) bool) OptionalFunc {
	return func(opt *optionalOptions) {
		opt.RequestIDSkipper = fn
	}
}

func WithMetricsSkipper(fn func(echo.Context) bool) OptionalFunc {
	return func(opt *optionalOptions) {
		opt.MetricsSkipper = fn
	}
}

func WithTraceSkipper(fn func(echo.Context) bool) OptionalFunc {
	return func(opt *optionalOptions) {
		opt.TraceSkipper = fn
	}
}

func WithLoggerSkipper(fn func(echo.Context) bool) OptionalFunc {
	return func(opt *optionalOptions) {
		opt.LoggerSkipper = fn
	}
}

func WithCustomMiddlewares(fns ...echo.MiddlewareFunc) OptionalFunc {
	return func(opt *optionalOptions) {
		opt.Middleware = fns
	}
}

func NewStrict(handlersRegistrator func(*echo.Echo), opt *Options, optionalFn ...OptionalFunc) (*Server, error) {
	return New(handlersRegistrator, opt, optionalFn...)
}

func New(handlersRegistrator func(*echo.Echo), opt *Options, optionalFn ...OptionalFunc) (*Server, error) {
	e := echo.New()
	e.Debug = opt.Debug
	e.HideBanner = true
	e.HidePort = true

	optionalOpts := optionalOptions{
		RequestIDSkipper: defaultURLSkipper,
		MetricsSkipper:   defaultURLSkipper,
		TraceSkipper:     defaultURLSkipper,
		LoggerSkipper:    defaultURLSkipper,
		Middleware:       nil,
	}

	for _, f := range optionalFn {
		f(&optionalOpts)
	}

	middleware.Use(e, &middleware.Options{
		ServiceName:      opt.ServiceName,
		Logger:           opt.Logger,
		TokenValidator:   opt.AuthService.ValidateToken,
		RequestIDSkipper: optionalOpts.RequestIDSkipper,
		MetricsSkipper:   optionalOpts.MetricsSkipper,
		TraceSkipper:     optionalOpts.TraceSkipper,
		LoggerSkipper:    optionalOpts.LoggerSkipper,
	})

	for _, mw := range optionalOpts.Middleware {
		e.Use(mw)
	}

	handlersRegistrator(e)

	return &Server{
		Echo:   e,
		Logger: opt.Logger,
	}, nil
}

type serveOptions struct {
	port                                                      int
	readTimeout, readHeaderTimeout, writeTimeout, idleTimeout time.Duration
}

type ServeOptFunc func(opt *serveOptions)

func WithPort(port int) ServeOptFunc {
	return func(opt *serveOptions) {
		opt.port = port
	}
}

func WithReadTimeout(timeout time.Duration) ServeOptFunc {
	return func(opt *serveOptions) {
		opt.readTimeout = timeout
	}
}

func WithReadHeaderTimeout(timeout time.Duration) ServeOptFunc {
	return func(opt *serveOptions) {
		opt.readHeaderTimeout = timeout
	}
}

func WithWriteTimeout(timeout time.Duration) ServeOptFunc {
	return func(opt *serveOptions) {
		opt.writeTimeout = timeout
	}
}

func WithIdleTimeout(timeout time.Duration) ServeOptFunc {
	return func(opt *serveOptions) {
		opt.idleTimeout = timeout
	}
}

func (s *Server) Serve(optionalFn ...ServeOptFunc) error {
	const (
		defaultServeTimeout      = 30 * time.Second
		defaultReadHeaderTimeout = 10 * time.Second
	)

	opts := serveOptions{
		port: 8000,

		readTimeout:       defaultServeTimeout,
		readHeaderTimeout: defaultReadHeaderTimeout,
		writeTimeout:      defaultServeTimeout,
		idleTimeout:       defaultServeTimeout,
	}

	for _, f := range optionalFn {
		f(&opts)
	}

	httpServer := &http.Server{ //nolint:exhaustruct
		Addr:              fmt.Sprintf(":%d", opts.port),
		Handler:           nil,
		TLSConfig:         nil,
		ReadTimeout:       opts.readTimeout,
		ReadHeaderTimeout: opts.readHeaderTimeout,
		WriteTimeout:      opts.writeTimeout,
		IdleTimeout:       opts.idleTimeout,
		MaxHeaderBytes:    0,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          log.New(logger.NewHTTPServerLogger(s.Logger), "", 0),
		BaseContext:       nil,
		ConnContext:       nil,
	}

	s.Echo.Server = httpServer

	s.Logger.Info("starting HTTP server", zap.Int("port", opts.port))

	if err := s.Echo.StartServer(httpServer); err != nil {
		return fmt.Errorf("error while start HTTP server | %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if err := s.Echo.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %w", err)
	}

	return nil
}

func defaultURLSkipper(c echo.Context) bool {
	for _, prefix := range skipURLPrefixList {
		if strings.HasPrefix(c.Path(), prefix) {
			return true
		}
	}

	return false
}
