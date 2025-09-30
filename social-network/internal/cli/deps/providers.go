package deps

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/samber/do"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/shaelmaar/otus-highload/social-network/internal/config"
	"github.com/shaelmaar/otus-highload/social-network/internal/domain"
	friendHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers/friend"
	loadTestHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers/loadtest"
	postHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers/post"
	userHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers/user"
	"github.com/shaelmaar/otus-highload/social-network/internal/metrics"
	"github.com/shaelmaar/otus-highload/social-network/internal/queries/pg"
	friendRepo "github.com/shaelmaar/otus-highload/social-network/internal/repository/friend"
	loadTestRepo "github.com/shaelmaar/otus-highload/social-network/internal/repository/loadtest"
	postRepo "github.com/shaelmaar/otus-highload/social-network/internal/repository/post"
	userRepo "github.com/shaelmaar/otus-highload/social-network/internal/repository/user"
	"github.com/shaelmaar/otus-highload/social-network/internal/service/auth"
	friendUseCases "github.com/shaelmaar/otus-highload/social-network/internal/usecase/friend"
	loadTestUseCases "github.com/shaelmaar/otus-highload/social-network/internal/usecase/loadtest"
	postUseCases "github.com/shaelmaar/otus-highload/social-network/internal/usecase/post"
	userUseCases "github.com/shaelmaar/otus-highload/social-network/internal/usecase/user"
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
			do.MustInvoke[*transaction.TxExecutor](i),
		)
	})

	do.Provide(i, func(i *do.Injector) (*friendUseCases.UseCases, error) {
		return friendUseCases.New(
			do.MustInvoke[domain.FriendRepository](i),
		)
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

	do.Provide(i, func(i *do.Injector) (*loadTestHandlers.Handlers, error) {
		return loadTestHandlers.New(
			do.MustInvoke[*loadTestUseCases.UseCases](i),
			do.MustInvoke[*zap.Logger](i),
		)
	})
}

func provideAuthService(i *do.Injector, cfg *config.Config) {
	do.Provide(i, func(i *do.Injector) (*auth.Service, error) {
		return auth.NewService(cfg.Auth.SecretKey, cfg.Auth.Expiration, cfg.ServiceName)
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
