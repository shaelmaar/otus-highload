package deps

import (
	"fmt"
	"net"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"github.com/tarantool/go-tarantool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/config"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/domain"
	httpHanders "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/httptransport/handlers"
	dialogHandlers "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/httptransport/handlers/dialog"
	httpServer "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/httptransport/server"
	dialogRepo "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/repository/dialog"
	dialogUseCases "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/usecase/dialog"
)

func provideUseCases(i *do.Injector) {
	do.Provide(i, func(i *do.Injector) (*dialogUseCases.UseCases, error) {
		return dialogUseCases.New(do.MustInvoke[domain.DialogRepository](i))
	})
}

func provideHTTPHandlers(i *do.Injector) {
	do.Provide(i, func(i *do.Injector) (*dialogHandlers.Handlers, error) {
		return dialogHandlers.New(
			do.MustInvoke[*dialogUseCases.UseCases](i),
			do.MustInvoke[*zap.Logger](i),
		)
	})

	do.ProvideNamed(i, nameHTTPHandlers, func(i *do.Injector) (*httpHanders.Handlers, error) {
		return httpHanders.New(do.MustInvoke[*dialogHandlers.Handlers](i))
	})
}

func provideHTTPServer(c *Container, cfg *config.Config) {
	do.ProvideNamed(c.i, nameHTTPServer, func(i *do.Injector) (*httpServer.Server, error) {
		s, err := httpServer.NewStrict(
			func(e *echo.Echo) {
				si := serverhttp.NewStrictHandler(
					do.MustInvokeNamed[*httpHanders.Handlers](i, nameHTTPHandlers), nil)

				serverhttp.RegisterHandlers(e, si)
			},

			&httpServer.Options{
				Debug:       false,
				ServiceName: cfg.ServiceName,
				Logger:      do.MustInvoke[*zap.Logger](i),
			},
		)
		if err != nil {
			return nil, fmt.Errorf("failed to init http server: %w", err)
		}

		c.addShutdown(nameHTTPServer, s.Stop)

		return s, nil
	})
}

func provideRepositories(i *do.Injector) {
	do.Provide(i, func(i *do.Injector) (domain.DialogRepository, error) {
		return dialogRepo.New(
			do.MustInvoke[*tarantool.Connection](i),
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

func provideTarantoolConnection(cfg *config.Config) (*tarantool.Connection, error) {
	conn, err := tarantool.Connect(net.JoinHostPort(cfg.TarantoolDB.Host, cfg.TarantoolDB.Port),
		//nolint:exhaustruct // остальное по умолчанию.
		tarantool.Opts{
			User: cfg.TarantoolDB.User,
			Pass: cfg.TarantoolDB.Pass,
		})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to TarantoolDB: %w", err)
	}

	return conn, nil
}
