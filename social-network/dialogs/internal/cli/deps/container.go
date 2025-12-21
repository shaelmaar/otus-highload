package deps

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/config"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/debugserver"
)

const (
	nameTarantoolConnection      = "tarantoolConnection"
	nameDebugServer              = "debugServer"
	nameHTTPServer               = "httpServer"
	nameHTTPHandlers             = "httpHandlers"
	nameDialogsMessagesProducer  = "dialogsMessagesProducer"
	nameCountersMessagesConsumer = "countersMessagesConsumer"
)

type shutdownFunc func(ctx context.Context) error

func sdWithoutCtx(f func() error) shutdownFunc {
	return func(_ context.Context) error {
		return f()
	}
}

type shutdown struct {
	sd   shutdownFunc
	name string
}

type Container struct {
	i *do.Injector

	shutdowns []shutdown
	mu        sync.Mutex
}

//nolint:funlen // инициализация DI контейнера.
func New(ctx context.Context) (*Container, error) {
	i := do.New()
	c := &Container{
		i:         i,
		shutdowns: nil,
		mu:        sync.Mutex{},
	}

	// low level deps, config, transports, clients etc
	do.Provide(i, func(i *do.Injector) (*config.Config, error) {
		return provideConfig()
	})

	cfg := do.MustInvoke[*config.Config](i)

	do.Provide(i, func(i *do.Injector) (*zap.Logger, error) {
		return provideLogger(cfg), nil
	})

	do.Provide(i, func(i *do.Injector) (*tarantool.Connection, error) {
		conn, err := provideTarantoolConnection(cfg)
		if err != nil {
			return nil, err
		}

		c.addShutdown(nameTarantoolConnection, sdWithoutCtx(conn.Close))

		return conn, nil
	})

	do.ProvideNamed(i, nameDebugServer, func(i *do.Injector) (*echo.Echo, error) {
		debugServer := debugserver.New()

		c.addShutdown(nameDebugServer, debugServer.Shutdown)

		return debugServer, nil
	})

	provideRepositories(i)

	provideKafkaProducers(c, cfg)

	provideKafkaConsumers(c, cfg)

	provideUseCases(i)

	provideHTTPHandlers(i)

	//nolint:contextcheck // контекст тут никак не передается.
	provideHTTPServer(c, cfg)

	return c, nil
}

// Shutdown - останавливает зависимости в порядке обратном инициализации.
func (c *Container) Shutdown(ctx context.Context) {
	logger := c.Logger()

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, sd := range slices.Backward(c.shutdowns) {
		logger.Info("shutting down " + sd.name)

		if err := sd.sd(ctx); err != nil {
			msg := fmt.Sprintf("error on %s shutdown", sd.name)
			logger.Warn(msg, zap.Error(err))
		}
	}
}

// addShutdown - регистрирует функцию останова для зависимости.
func (c *Container) addShutdown(name string, sd shutdownFunc) {
	c.mu.Lock()

	c.shutdowns = append(c.shutdowns, shutdown{name: name, sd: sd})

	c.mu.Unlock()
}
