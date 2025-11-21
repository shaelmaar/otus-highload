package deps

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"github.com/samber/do"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/shaelmaar/otus-highload/social-network/gen/serverhttp"
	postfeedCache "github.com/shaelmaar/otus-highload/social-network/internal/cache/postfeed"
	"github.com/shaelmaar/otus-highload/social-network/internal/config"
	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	"github.com/shaelmaar/otus-highload/social-network/internal/dto"
	dialogsGRPCClient "github.com/shaelmaar/otus-highload/social-network/internal/grpctransport/clients/dialogs"
	grpcHandlers "github.com/shaelmaar/otus-highload/social-network/internal/grpctransport/handlers"
	grpcServer "github.com/shaelmaar/otus-highload/social-network/internal/grpctransport/server"
	httpHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers"
	dialogHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers/dialog"
	friendHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers/friend"
	loadTestHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers/loadtest"
	postHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers/post"
	userHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers/user"
	wsHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers/ws"
	httpServer "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/server"
	"github.com/shaelmaar/otus-highload/social-network/internal/metrics"
	"github.com/shaelmaar/otus-highload/social-network/internal/queries/pg"
	"github.com/shaelmaar/otus-highload/social-network/internal/rabbitmq"
	friendRepo "github.com/shaelmaar/otus-highload/social-network/internal/repository/friend"
	loadTestRepo "github.com/shaelmaar/otus-highload/social-network/internal/repository/loadtest"
	postRepo "github.com/shaelmaar/otus-highload/social-network/internal/repository/post"
	userRepo "github.com/shaelmaar/otus-highload/social-network/internal/repository/user"
	"github.com/shaelmaar/otus-highload/social-network/internal/service/auth"
	"github.com/shaelmaar/otus-highload/social-network/internal/service/postfeed"
	"github.com/shaelmaar/otus-highload/social-network/internal/taskhandler/postcreatedchunked"
	"github.com/shaelmaar/otus-highload/social-network/internal/taskhandler/userupdatefeed"
	"github.com/shaelmaar/otus-highload/social-network/internal/taskhandler/userupdatefeedchunked"
	dialogUseCases "github.com/shaelmaar/otus-highload/social-network/internal/usecase/dialog"
	feedUseCases "github.com/shaelmaar/otus-highload/social-network/internal/usecase/feed"
	friendUseCases "github.com/shaelmaar/otus-highload/social-network/internal/usecase/friend"
	loadTestUseCases "github.com/shaelmaar/otus-highload/social-network/internal/usecase/loadtest"
	postUseCases "github.com/shaelmaar/otus-highload/social-network/internal/usecase/post"
	userUseCases "github.com/shaelmaar/otus-highload/social-network/internal/usecase/user"
	"github.com/shaelmaar/otus-highload/social-network/internal/valkeyprovider"
	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
)

func provideUseCases(i *do.Injector) {
	do.Provide(i, func(i *do.Injector) (*userUseCases.UseCases, error) {
		return userUseCases.New(
			do.MustInvoke[domain.UserRepository](i),
			do.MustInvoke[*auth.Service](i),
			do.MustInvoke[*transaction.TxExecutor](i),
		)
	})

	do.Provide(i, func(i *do.Injector) (*postUseCases.UseCases, error) {
		return postUseCases.New(
			do.MustInvoke[domain.PostRepository](i),
			do.MustInvoke[domain.FriendRepository](i),
			do.MustInvoke[*postfeed.Service](i),
			do.MustInvokeNamed[*rabbitmq.Producer[dto.UserFeedChunkedUpdateTask]](i, nameUserFeedChunkedTaskProducer),
			do.MustInvokeNamed[*rabbitmq.Producer[dto.PostCreatedChunkedTask]](i, namePostCreatedChunkedTaskProducer),
			do.MustInvoke[*transaction.TxExecutor](i),
		)
	})

	do.Provide(i, func(i *do.Injector) (*friendUseCases.UseCases, error) {
		return friendUseCases.New(
			do.MustInvoke[domain.FriendRepository](i),
			do.MustInvokeNamed[*rabbitmq.Producer[dto.UserFeedUpdateTask]](i, nameUserFeedTaskProducer),
		)
	})

	do.Provide(i, func(i *do.Injector) (*feedUseCases.UseCases, error) {
		return feedUseCases.New(do.MustInvoke[*postfeed.Service](i))
	})

	do.Provide(i, func(i *do.Injector) (*dialogUseCases.UseCases, error) {
		return dialogUseCases.New(
			do.MustInvokeNamed[*dialogsGRPCClient.Client](i, nameDialogsGRPCClient).DialogsService)
	})

	do.Provide(i, func(i *do.Injector) (*loadTestUseCases.UseCases, error) {
		return loadTestUseCases.New(
			do.MustInvoke[domain.LoadTestRepository](i),
			do.MustInvoke[*transaction.TxExecutor](i),
			do.MustInvoke[*metrics.Metrics](i),
			do.MustInvoke[*zap.Logger](i),
		)
	})
}

func provideHTTPHandlers(i *do.Injector) {
	do.Provide(i, func(i *do.Injector) (*userHandlers.Handlers, error) {
		return userHandlers.New(
			do.MustInvoke[*userUseCases.UseCases](i),
			do.MustInvoke[*zap.Logger](i),
		)
	})

	do.Provide(i, func(i *do.Injector) (*postHandlers.Handlers, error) {
		return postHandlers.New(
			do.MustInvoke[*postUseCases.UseCases](i),
			do.MustInvoke[*zap.Logger](i),
		)
	})

	do.Provide(i, func(i *do.Injector) (*friendHandlers.Handlers, error) {
		return friendHandlers.New(
			do.MustInvoke[*friendUseCases.UseCases](i),
			do.MustInvoke[*zap.Logger](i),
		)
	})

	do.Provide(i, func(i *do.Injector) (*dialogHandlers.Handlers, error) {
		return dialogHandlers.New(
			do.MustInvoke[*dialogUseCases.UseCases](i),
			do.MustInvoke[*zap.Logger](i),
		)
	})

	do.Provide(i, func(i *do.Injector) (*loadTestHandlers.Handlers, error) {
		return loadTestHandlers.New(
			do.MustInvoke[*loadTestUseCases.UseCases](i),
			do.MustInvoke[*zap.Logger](i),
		)
	})

	do.ProvideNamed[*httpHandlers.Handlers](i, nameHTTPHandlers,
		func(i *do.Injector) (*httpHandlers.Handlers, error) {
			return httpHandlers.NewHandlers(
				do.MustInvoke[*userHandlers.Handlers](i),
				do.MustInvoke[*postHandlers.Handlers](i),
				do.MustInvoke[*friendHandlers.Handlers](i),
				do.MustInvoke[*dialogHandlers.Handlers](i),
				do.MustInvoke[*loadTestHandlers.Handlers](i),
			)
		})
}

func provideMelody(c *Container) {
	do.Provide(c.i, func(i *do.Injector) (*melody.Melody, error) {
		m := melody.New()

		// конфигурацию можно изменить в m.Config.
		// пока дефолтного конфига достаточно.

		c.addShutdown(nameMelodyWebsocket, sdWithoutCtx(m.Close))

		return m, nil
	})
}

func provideWSServer(c *Container, cfg *config.Config) {
	do.Provide(c.i, func(i *do.Injector) (*wsHandlers.Handlers, error) {
		return wsHandlers.New(
			do.MustInvoke[*melody.Melody](i),
			do.MustInvoke[*zap.Logger](i),
		)
	})

	do.Provide(c.i, func(i *do.Injector) (*httpHandlers.WSHandlers, error) {
		return httpHandlers.NewWSHandlers(
			do.MustInvoke[*wsHandlers.Handlers](i),
		)
	})

	do.ProvideNamed(c.i, nameWSServer, func(i *do.Injector) (*httpServer.Server, error) {
		s, err := httpServer.New(
			httpServer.RegisterWSHandlers(do.MustInvoke[*httpHandlers.WSHandlers](i)),
			&httpServer.Options{
				Debug:       false,
				ServiceName: cfg.ServiceName,
				Logger:      do.MustInvoke[*zap.Logger](i),
				AuthService: do.MustInvoke[*auth.Service](i),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to init websocket server: %w", err)
		}

		c.addShutdown(nameWSServer, s.Stop)

		return s, nil
	})
}

func provideHTTPServer(c *Container, cfg *config.Config) {
	do.ProvideNamed(c.i, nameHTTPServer, func(i *do.Injector) (*httpServer.Server, error) {
		logger := do.MustInvoke[*zap.Logger](i)

		s, err := httpServer.NewStrict(
			func(e *echo.Echo) {
				si := serverhttp.NewStrictHandler(
					do.MustInvokeNamed[*httpHandlers.Handlers](i, nameHTTPHandlers),
					nil)

				serverhttp.RegisterHandlers(e, si)
			},

			&httpServer.Options{
				Debug:       false,
				ServiceName: cfg.ServiceName,
				Logger:      logger,
				AuthService: do.MustInvoke[*auth.Service](i),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to init http server: %w", err)
		}

		c.addShutdown(nameHTTPServer, s.Stop)

		return s, nil
	})
}

func provideGRPCHandlers(i *do.Injector) {
	do.ProvideNamed(i, nameGRPCHandlers, func(i *do.Injector) (*grpcHandlers.Handlers, error) {
		return grpcHandlers.New(do.MustInvoke[*auth.Service](i))
	})
}

func provideGRPCServer(c *Container) {
	do.ProvideNamed(c.i, nameGRPCServer, func(i *do.Injector) (*grpcServer.Server, error) {
		s, err := grpcServer.New(&grpcServer.NewServerOptions{
			Logger:            do.MustInvoke[*zap.Logger](i),
			GRPCHandlers:      do.MustInvokeNamed[*grpcHandlers.Handlers](i, nameGRPCHandlers),
			Validator:         nil,
			UnaryInterceptors: nil,
			ServerOptions:     nil,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to init grpc server: %w", err)
		}

		c.addShutdown(nameGRPCServer, sdSimple(s.Stop))

		return s, nil
	})
}

func provideAuthService(i *do.Injector, cfg *config.Config) {
	do.Provide(i, func(i *do.Injector) (*auth.Service, error) {
		return auth.NewService(cfg.Auth.SecretKey, cfg.Auth.Expiration, cfg.ServiceName)
	})
}

func providePostFeedService(i *do.Injector) {
	do.Provide(i, func(injector *do.Injector) (*postfeed.Service, error) {
		return postfeed.NewService(
			do.MustInvoke[domain.PostRepository](i),
			do.MustInvoke[domain.FriendRepository](i),
			do.MustInvoke[*postfeedCache.Cache](i),
			do.MustInvoke[*zap.Logger](i),
		)
	})
}

func provideCaches(i *do.Injector) {
	do.Provide(i, func(injector *do.Injector) (*postfeedCache.Cache, error) {
		return postfeedCache.New(
			do.MustInvoke[*valkeyprovider.Provider](i),
			do.MustInvoke[*zap.Logger](i),
		)
	})
}

func provideRepositories(i *do.Injector) {
	do.Provide(i, func(i *do.Injector) (domain.UserRepository, error) {
		return userRepo.New(
			do.MustInvokeNamed[pg.QuerierTX](i, nameQuerier),
			do.MustInvokeNamed[pg.QuerierTX](i, nameReplicaQuerier),
		)
	})

	do.Provide(i, func(i *do.Injector) (domain.PostRepository, error) {
		return postRepo.New(
			do.MustInvokeNamed[pg.QuerierTX](i, nameQuerier),
			do.MustInvokeNamed[pg.QuerierTX](i, nameReplicaQuerier),
		)
	})

	do.Provide(i, func(i *do.Injector) (domain.FriendRepository, error) {
		return friendRepo.New(
			do.MustInvokeNamed[pg.QuerierTX](i, nameQuerier),
			do.MustInvokeNamed[pg.QuerierTX](i, nameReplicaQuerier),
		)
	})

	do.Provide(i, func(i *do.Injector) (domain.LoadTestRepository, error) {
		return loadTestRepo.New(
			do.MustInvokeNamed[pg.QuerierTX](i, nameQuerier),
		)
	})
}

func provideTaskProducers(c *Container, cfg *config.Config) {
	do.ProvideNamed(
		c.i, nameUserFeedTaskProducer, func(i *do.Injector) (*rabbitmq.Producer[dto.UserFeedUpdateTask], error) {
			p, err := rabbitmq.NewProducer[dto.UserFeedUpdateTask](
				cfg.RabbitMQ.URL(),
				rabbitmq.UserFeedUpdateQueueName,
				"",
				do.MustInvoke[*zap.Logger](i),
			)
			if err != nil {
				return nil, err
			}

			c.addShutdown(nameUserFeedTaskProducer, sdSimple(p.Close))

			return p, nil
		},
	)

	do.ProvideNamed(
		c.i, nameUserFeedChunkedTaskProducer,
		func(i *do.Injector) (*rabbitmq.Producer[dto.UserFeedChunkedUpdateTask], error) {
			p, err := rabbitmq.NewProducer[dto.UserFeedChunkedUpdateTask](
				cfg.RabbitMQ.URL(),
				rabbitmq.UserFeedUpdateChunkedQueueName,
				"",
				do.MustInvoke[*zap.Logger](i),
			)
			if err != nil {
				return nil, err
			}

			c.addShutdown(nameUserFeedChunkedTaskProducer, sdSimple(p.Close))

			return p, nil
		},
	)

	do.ProvideNamed(
		c.i, namePostCreatedChunkedTaskProducer,
		func(i *do.Injector) (*rabbitmq.Producer[dto.PostCreatedChunkedTask], error) {
			p, err := rabbitmq.NewProducer[dto.PostCreatedChunkedTask](
				cfg.RabbitMQ.URL(),
				"",
				rabbitmq.PostCreatedExchangeName,
				do.MustInvoke[*zap.Logger](i),
			)
			if err != nil {
				return nil, err
			}

			c.addShutdown(namePostCreatedChunkedTaskProducer, sdSimple(p.Close))

			return p, nil
		},
	)
}

func provideTaskConsumers(c *Container, cfg *config.Config) {
	do.Provide(c.i, func(i *do.Injector) (*userupdatefeed.Handler, error) {
		return userupdatefeed.New(
			do.MustInvoke[*feedUseCases.UseCases](i),
		)
	})

	do.ProvideNamed(
		c.i, nameUserFeedTaskConsumer, func(i *do.Injector) (*rabbitmq.Consumer[dto.UserFeedUpdateTask], error) {
			consumer, err := rabbitmq.NewConsumer[dto.UserFeedUpdateTask](
				cfg.RabbitMQ.URL(),
				rabbitmq.UserFeedUpdateQueueName,
				"",
				do.MustInvoke[*userupdatefeed.Handler](i).Handle,
				do.MustInvoke[*zap.Logger](i),
			)
			if err != nil {
				return nil, err
			}

			c.addShutdown(nameUserFeedTaskConsumer, sdSimple(consumer.Close))

			return consumer, nil
		},
	)

	do.Provide(c.i, func(i *do.Injector) (*userupdatefeedchunked.Handler, error) {
		return userupdatefeedchunked.New(
			do.MustInvoke[*feedUseCases.UseCases](i),
		)
	})

	do.ProvideNamed(
		c.i, nameUserFeedChunkedTaskConsumer,
		func(i *do.Injector) (*rabbitmq.Consumer[dto.UserFeedChunkedUpdateTask], error) {
			consumer, err := rabbitmq.NewConsumer[dto.UserFeedChunkedUpdateTask](
				cfg.RabbitMQ.URL(),
				rabbitmq.UserFeedUpdateChunkedQueueName,
				"",
				do.MustInvoke[*userupdatefeedchunked.Handler](i).Handle,
				do.MustInvoke[*zap.Logger](i),
			)
			if err != nil {
				return nil, err
			}

			c.addShutdown(nameUserFeedChunkedTaskConsumer, sdSimple(consumer.Close))

			return consumer, nil
		},
	)

	do.Provide(c.i, func(i *do.Injector) (*postcreatedchunked.Handler, error) {
		return postcreatedchunked.New(
			do.MustInvoke[*melody.Melody](i),
		)
	})

	do.ProvideNamed(
		c.i, namePostCreatedChunkedTaskConsumer,
		func(i *do.Injector) (*rabbitmq.Consumer[dto.PostCreatedChunkedTask], error) {
			consumer, err := rabbitmq.NewConsumer[dto.PostCreatedChunkedTask](
				cfg.RabbitMQ.URL(),
				"",
				rabbitmq.PostCreatedExchangeName,
				do.MustInvoke[*postcreatedchunked.Handler](i).Handle,
				do.MustInvoke[*zap.Logger](i),
			)
			if err != nil {
				return nil, err
			}

			c.addShutdown(namePostCreatedChunkedTaskConsumer, sdSimple(consumer.Close))

			return consumer, nil
		},
	)
}

func provideGRPCClients(i *do.Injector, cfg *config.Config) {
	do.ProvideNamed[*dialogsGRPCClient.Client](
		i, nameDialogsGRPCClient, func(i *do.Injector) (*dialogsGRPCClient.Client, error) {
			c, err := dialogsGRPCClient.NewGRPCClient(&dialogsGRPCClient.NewClientOptions{
				GRPCAddr:          cfg.DialogsGRPCClient.Host,
				TLS:               cfg.DialogsGRPCClient.TLS,
				Timeout:           &cfg.DialogsGRPCClient.Timeout,
				UnaryInterceptors: nil,
				DialOptions:       nil,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to init dialogs grpc client: %w", err)
			}

			return c, nil
		})
}

func provideConfig() (*config.Config, error) {
	cfg, err := config.FromEnv()
	if err != nil {
		return nil, err //nolint:wrapcheck
	}

	return cfg, nil
}

func provideLogger(cfg *config.Config) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewJSONEncoder(encoderCfg)

	stacktraceLevel := zapcore.ErrorLevel
	if !cfg.Log.EnableStacktrace {
		stacktraceLevel = zapcore.FatalLevel + 1
	}

	stdoutFilter := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		if cfg.Debug {
			return level < zapcore.ErrorLevel
		}

		return level > zapcore.DebugLevel && level < zapcore.ErrorLevel
	})

	stderrFilter := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stdout),
			stdoutFilter,
		),
		zapcore.NewCore(
			encoder,
			zapcore.Lock(os.Stderr),
			stderrFilter,
		),
	)

	return zap.New(core, zap.AddStacktrace(stacktraceLevel))
}

func providePostgresql(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	pgxCfg, err := cfg.Database.PgxConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgreSQL: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping postgreSQL: %w", err)
	}

	return pool, nil
}

func provideReplicaPostgresql(ctx context.Context, cfg *config.Config, pgxPool *pgxpool.Pool) (*pgxpool.Pool, error) {
	if !cfg.ReplicaDatabase.Enabled {
		return pgxPool, nil
	}

	pgxCfg, err := cfg.ReplicaDatabase.PgxConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgreSQL: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping postgreSQL: %w", err)
	}

	return pool, nil
}
