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
	monolithGPRCClient "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/grpctransport/clients/monolith"
	grpcHandlers "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/grpctransport/handlers"
	grpcServer "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/grpctransport/server"
	httpHanders "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/httptransport/handlers"
	dialogHandlers "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/httptransport/handlers/dialog"
	httpServer "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/httptransport/server"
	dialogRepo "github.com/shaelmaar/otus-highload/social-network/dialogs/internal/repository/dialog"
	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/service/auth"
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
		return grpcHandlers.New(do.MustInvoke[*dialogUseCases.UseCases](i))
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

func provideGRPCClients(i *do.Injector, cfg *config.Config) {
	do.ProvideNamed[*monolithGPRCClient.Client](
		i, nameMonolithGRPCClient, func(i *do.Injector) (*monolithGPRCClient.Client, error) {
			c, err := monolithGPRCClient.NewGRPCClient(&monolithGPRCClient.NewClientOptions{
				GRPCAddr:          cfg.MonolithGRPCClient.Host,
				TLS:               cfg.MonolithGRPCClient.TLS,
				Timeout:           &cfg.MonolithGRPCClient.Timeout,
				UnaryInterceptors: nil,
				DialOptions:       nil,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to init monolith grpc client: %w", err)
			}

			return c, nil
		})
}

func provideAuthService(i *do.Injector) {
	do.Provide(i, func(i *do.Injector) (*auth.Service, error) {
		return auth.New(do.MustInvokeNamed[*monolithGPRCClient.Client](i, nameMonolithGRPCClient).AuthService)
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
