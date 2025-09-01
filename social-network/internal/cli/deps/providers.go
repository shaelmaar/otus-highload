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
	loadTestHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers/loadtest"
	userHandlers "github.com/shaelmaar/otus-highload/social-network/internal/httptransport/handlers/user"
	loadTestUseCases "github.com/shaelmaar/otus-highload/social-network/internal/usecase/loadtest"
	userUseCases "github.com/shaelmaar/otus-highload/social-network/internal/usecase/user"
	"github.com/shaelmaar/otus-highload/social-network/pkg/transaction"
)

func provideHTTPHandlers(i *do.Injector) {
	do.Provide(i, func(i *do.Injector) (*userUseCases.UseCases, error) {
		return userUseCases.New(
			do.MustInvoke[domain.UserRepository](i),
			do.MustInvoke[*transaction.TxExecutor](i),
		)
	})

	do.Provide(i, func(i *do.Injector) (*loadTestUseCases.UseCases, error) {
		return loadTestUseCases.New(
			do.MustInvoke[domain.LoadTestRepository](i),
		)
	})

	do.Provide(i, func(i *do.Injector) (*userHandlers.Handlers, error) {
		return userHandlers.NewHandlers(
			do.MustInvoke[*userUseCases.UseCases](i),
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
