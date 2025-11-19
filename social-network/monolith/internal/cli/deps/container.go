package deps

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"github.com/tarantool/go-tarantool"
	"github.com/valkey-io/valkey-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/config"
	"github.com/shaelmaar/otus-highload/social-network/internal/debugserver"
	"github.com/shaelmaar/otus-highload/social-network/internal/metrics"
	"github.com/shaelmaar/otus-highload/social-network/internal/queries/pg"
	"github.com/shaelmaar/otus-highload/social-network/internal/valkeyprovider"
	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
)

const (
	namePgxPool                        = "pgxPool"
	nameReplicaPgxPool                 = "replicaPgxPool"
	nameMongoDialogsDB                 = "mongoDialogsDB"
	nameTarantoolConnection            = "tarantoolConnection"
	nameQuerier                        = "querier"
	nameReplicaQuerier                 = "replicaQuerier"
	nameDebugServer                    = "debugServer"
	nameValkeyProvider                 = "valkeyProvider"
	nameUserFeedTaskProducer           = "userFeedTaskProducer"
	nameUserFeedChunkedTaskProducer    = "userFeedChunkedTaskProducer"
	namePostCreatedChunkedTaskProducer = "postCreatedChunkedTaskProducer"
	nameUserFeedTaskConsumer           = "userFeedTaskConsumer"
	nameUserFeedChunkedTaskConsumer    = "userFeedChunkedTaskConsumer"
	namePostCreatedChunkedTaskConsumer = "postCreatedChunkedTaskConsumer"
	nameMelodyWebsocket                = "melodyWebsocket"
	nameWSServer                       = "wsServer"
	nameGRPCServer                     = "grpcServer"
	nameHTTPServer                     = "httpServer"
)

type shutdownFunc func(ctx context.Context) error

func sdSimple(f func()) shutdownFunc {
	return func(ctx context.Context) error {
		f()

		return nil
	}
}

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

	do.ProvideNamed(i, namePgxPool, func(i *do.Injector) (*pgxpool.Pool, error) {
		pool, err := providePostgresql(ctx, cfg)
		if err != nil {
			return nil, err
		}

		c.addShutdown(namePgxPool, sdSimple(pool.Close))

		return pool, nil
	})

	do.ProvideNamed(i, nameReplicaPgxPool, func(i *do.Injector) (*pgxpool.Pool, error) {
		pool, err := provideReplicaPostgresql(
			ctx, cfg, do.MustInvokeNamed[*pgxpool.Pool](i, namePgxPool),
		)
		if err != nil {
			return nil, err
		}

		c.addShutdown(nameReplicaPgxPool, sdSimple(pool.Close))

		return pool, nil
	})

	do.ProvideNamed(i, nameMongoDialogsDB, func(i *do.Injector) (*mongo.Database, error) {
		db, err := provideMongoDialogsDB(ctx, cfg)
		if err != nil {
			return nil, err
		}

		c.addShutdown(nameMongoDialogsDB, db.Client().Disconnect)

		return db, nil
	})

	do.Provide(i, func(i *do.Injector) (*tarantool.Connection, error) {
		conn, err := provideTarantoolConnection(cfg)
		if err != nil {
			return nil, err
		}

		c.addShutdown(nameTarantoolConnection, sdWithoutCtx(conn.Close))

		return conn, nil
	})

	do.ProvideNamed(i, nameQuerier, func(i *do.Injector) (pg.QuerierTX, error) {
		return pg.NewQueriesTX(pg.New(do.MustInvokeNamed[*pgxpool.Pool](i, namePgxPool))), nil
	})

	do.ProvideNamed(i, nameReplicaQuerier, func(i *do.Injector) (pg.QuerierTX, error) {
		return pg.NewQueriesTX(pg.New(do.MustInvokeNamed[*pgxpool.Pool](i, nameReplicaPgxPool))), nil
	})

	do.Provide(i, func(i *do.Injector) (*valkeyprovider.Provider, error) {
		p, err := valkeyprovider.NewProvider(
			//nolint:exhaustruct // остальное по умолчанию.
			valkey.ClientOption{
				InitAddress:      []string{cfg.Valkey.Address},
				SelectDB:         cfg.Valkey.DB,
				ConnWriteTimeout: cfg.Valkey.SetTimeout,
				Password:         cfg.Valkey.Password,
			},
			do.MustInvoke[*zap.Logger](i),
		)
		if err != nil {
			return nil, err
		}

		p.Init()

		c.addShutdown(nameValkeyProvider, sdSimple(p.Close))

		return p, nil
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

	provideCaches(i)

	provideTaskProducers(c, cfg)

	provideTaskConsumers(c, cfg)

	provideAuthService(i, cfg)

	providePostFeedService(i)

	provideUseCases(i)

	provideHTTPHandlers(i)

	provideMelody(c)

	//nolint:contextcheck // контекст тут никак не передается.
	provideWSServer(c, cfg)

	//nolint:contextcheck // контекст тут никак не передается.
	provideHTTPServer(c, cfg)

	provideGRPCHandlers(i)

	provideGRPCServer(c)

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
