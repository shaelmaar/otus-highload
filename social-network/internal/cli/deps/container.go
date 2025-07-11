package deps

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/config"
	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers"
	userHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers/user"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/server"
	"github.com/shaelmaar/otus-highload/social-network/internal/queries/pg"
	userRepo "github.com/shaelmaar/otus-highload/social-network/internal/repository/user"
	userUseCases "github.com/shaelmaar/otus-highload/social-network/internal/usecase/user"
	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
)

const (
	namePgxPool = "pgxPool"
	nameQuerier = "querier"
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

	do.ProvideNamed(i, nameQuerier, func(i *do.Injector) (pg.QuerierTX, error) {
		return pg.NewQueriesTX(pg.New(do.MustInvokeNamed[*pgxpool.Pool](i, namePgxPool))), nil
	})

	do.Provide(i, func(i *do.Injector) (*transaction.TxExecutor, error) {
		return transaction.New(
			do.MustInvokeNamed[*pgxpool.Pool](i, namePgxPool),
		)
	})

	do.Provide(i, func(i *do.Injector) (domain.UserRepository, error) {
		return userRepo.New(
			do.MustInvokeNamed[pg.QuerierTX](i, nameQuerier))
	})

	do.Provide(i, func(i *do.Injector) (*userUseCases.UseCases, error) {
		return userUseCases.New(
			do.MustInvoke[domain.UserRepository](i),
			do.MustInvoke[*transaction.TxExecutor](i),
		)
	})

	do.Provide[*userHandlers.Handlers](i, func(i *do.Injector) (*userHandlers.Handlers, error) {
		return userHandlers.NewHandlers(
			do.MustInvoke[*userUseCases.UseCases](i),
			do.MustInvoke[*zap.Logger](i),
		)
	})

	//nolint:contextcheck // контекст тут никак не передается.
	do.Provide(i, func(i *do.Injector) (*server.Server, error) {
		logger := do.MustInvoke[*zap.Logger](i)

		httpServer, err := server.NewStrict(
			handlers.NewHandlers(
				do.MustInvoke[*userHandlers.Handlers](i),
			),
			&server.Options{
				Debug:       false,
				ServiceName: cfg.ServiceName,
				Logger:      logger,
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
