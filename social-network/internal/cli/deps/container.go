package deps

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/config"
	"github.com/shaelmaar/otus-highload/social-network/internal/debugserver"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers"
	loadTestHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers/loadtest"
	postHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers/post"
	userHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers/user"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/server"
	"github.com/shaelmaar/otus-highload/social-network/internal/metrics"
	"github.com/shaelmaar/otus-highload/social-network/internal/queries/pg"
	"github.com/shaelmaar/otus-highload/social-network/internal/service/auth"
	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
)

const (
	namePgxPool        = "pgxPool"
	nameReplicaPgxPool = "replicaPgxPool"
	nameQuerier        = "querier"
	nameReplicaQuerier = "replicaQuerier"
	nameDebugServer    = "debugServer"
)

type shutdownFunc func(ctx context.Context) error

type shutdown struct {
	sd   shutdownFunc
	name string
}

type Container struct {
	i *do.Injector

	shutdowns []shutdown
	mu        sync.Mutex
}

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

	do.ProvideNamed(i, namePgxPool, func(i *do.Injector) (*pgxpool.Pool, error) {
		return providePostgresql(ctx, cfg)
	})

	do.ProvideNamed(i, nameReplicaPgxPool, func(i *do.Injector) (*pgxpool.Pool, error) {
		return provideReplicaPostgresql(
			ctx, cfg, do.MustInvokeNamed[*pgxpool.Pool](i, namePgxPool),
		)
	})

	do.ProvideNamed(i, nameQuerier, func(i *do.Injector) (pg.QuerierTX, error) {
		return pg.NewQueriesTX(pg.New(do.MustInvokeNamed[*pgxpool.Pool](i, namePgxPool))), nil
	})

	do.ProvideNamed(i, nameReplicaQuerier, func(i *do.Injector) (pg.QuerierTX, error) {
		return pg.NewQueriesTX(pg.New(do.MustInvokeNamed[*pgxpool.Pool](i, nameReplicaPgxPool))), nil
	})

	do.Provide(i, func(i *do.Injector) (*transaction.TxExecutor, error) {
		return transaction.New(
			do.MustInvokeNamed[*pgxpool.Pool](i, namePgxPool),
		)
	})

	do.Provide(i, func(i *do.Injector) (*metrics.Metrics, error) {
		return metrics.NewMetrics(), nil
	})

	do.ProvideNamed(i, nameDebugServer, func(i *do.Injector) (*echo.Echo, error) {
		debugServer := debugserver.New()

		c.addShutdown(nameDebugServer, debugServer.Shutdown)

		return debugServer, nil
	})

	provideRepositories(i)

	provideAuthService(i, cfg)

	provideUseCases(i)

	provideHTTPHandlers(i)

	//nolint:contextcheck // контекст тут никак не передается.
	do.Provide(i, func(i *do.Injector) (*server.Server, error) {
		logger := do.MustInvoke[*zap.Logger](i)

		httpServer, err := server.NewStrict(
			handlers.NewHandlers(
				do.MustInvoke[*userHandlers.Handlers](i),
				do.MustInvoke[*postHandlers.Handlers](i),
				do.MustInvoke[*loadTestHandlers.Handlers](i),
			),
			&server.Options{
				Debug:       false,
				ServiceName: cfg.ServiceName,
				Logger:      logger,
				AuthService: do.MustInvoke[*auth.Service](i),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to init http server: %w", err)
		}

		c.addShutdown("httpServer", httpServer.Stop)

		return httpServer, nil
	})

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
