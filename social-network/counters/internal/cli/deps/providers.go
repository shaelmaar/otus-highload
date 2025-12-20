package deps

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"github.com/segmentio/kafka-go"
	"github.com/valkey-io/valkey-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/shaelmaar/otus-highload/social-network/counters/gen/serverhttp"
	"github.com/shaelmaar/otus-highload/social-network/counters/internal/config"
	"github.com/shaelmaar/otus-highload/social-network/counters/internal/domain"
	httpHanders "github.com/shaelmaar/otus-highload/social-network/counters/internal/httptransport/handlers"
	dialogHandlers "github.com/shaelmaar/otus-highload/social-network/counters/internal/httptransport/handlers/dialog"
	httpServer "github.com/shaelmaar/otus-highload/social-network/counters/internal/httptransport/server"
	"github.com/shaelmaar/otus-highload/social-network/counters/internal/kafka/consumer/dialogsmessages"
	"github.com/shaelmaar/otus-highload/social-network/counters/internal/kafka/producer/messages"
	dialogRepo "github.com/shaelmaar/otus-highload/social-network/counters/internal/repository/dialog"
	dialogUseCases "github.com/shaelmaar/otus-highload/social-network/counters/internal/usecase/dialog"
)

func provideUseCases(i *do.Injector) {
	do.Provide(i, func(i *do.Injector) (*dialogUseCases.UseCases, error) {
		return dialogUseCases.New(
			do.MustInvoke[domain.DialogRepository](i),
			do.MustInvoke[*messages.Producer](i),
		)
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
			do.MustInvokeNamed[valkey.Client](i, nameValkeyMasterClient),
			do.MustInvokeNamed[[]valkey.Client](i, nameValkeyReplicaClients),
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

func provideValkeyClients(c *Container, cfg *config.Config) {
	do.ProvideNamed(c.i, nameValkeyMasterClient, func(i *do.Injector) (valkey.Client, error) {
		client, err := valkey.NewClient(valkey.ClientOption{
			InitAddress:      cfg.ValkeyDB.Addresses,
			ConnWriteTimeout: cfg.ValkeyDB.WriteTimeout,
			Password:         cfg.ValkeyDB.Password,
			Sentinel: valkey.SentinelOption{
				Dialer: net.Dialer{
					Timeout: time.Second,
				},
				MasterSet: "mymaster",
				Password:  cfg.ValkeyDB.Password,
			},
		})
		if err != nil {
			return nil, err
		}

		c.addShutdown(nameValkeyMasterClient, sdSimple(client.Close))

		return client, nil
	})

	do.ProvideNamed(c.i, nameValkeyReplicaClients, func(i *do.Injector) ([]valkey.Client, error) {
		clients := make([]valkey.Client, 0, cfg.ValkeyDB.SlavePoolSize)

		for i := 0; i < cfg.ValkeyDB.SlavePoolSize; i++ {
			client, err := valkey.NewClient(valkey.ClientOption{
				InitAddress: cfg.ValkeyDB.Addresses,
				ReplicaOnly: true,
				Password:    cfg.ValkeyDB.Password,
				Sentinel: valkey.SentinelOption{
					Dialer: net.Dialer{
						Timeout: time.Second,
					},
					MasterSet: "mymaster",
					Password:  cfg.ValkeyDB.Password,
				},
				ReadNodeSelector: func(_ uint16, nodes []valkey.NodeInfo) int {
					return rand.Intn(10)
				},
			})
			if err != nil {
				return nil, err
			}

			clients = append(clients, client)
		}

		c.addShutdown(nameValkeyReplicaClients, func(_ context.Context) error {
			for i := 0; i < cfg.ValkeyDB.SlavePoolSize; i++ {
				clients[i].Close()
			}

			return nil
		})

		return clients, nil
	})
}

func provideKafkaProducers(c *Container, cfg *config.Config) {
	do.Provide(c.i, func(i *do.Injector) (*messages.Producer, error) {
		kafkaWriter := &kafka.Writer{
			Addr:         kafka.TCP(cfg.Kafka.Brokers...),
			BatchSize:    1,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireOne,
			Topic:        messages.TopicName,
		}

		producer, err := messages.New(kafkaWriter)
		if err != nil {
			return nil, err
		}

		c.addShutdown(nameMessagesProducer, sdWithoutCtx(producer.Close))

		return producer, nil
	})
}

func provideKafkaConsumers(c *Container, cfg *config.Config) {
	do.ProvideNamed(c.i, nameDialogsMessagesConsumer, func(i *do.Injector) (*dialogsmessages.Consumer, error) {
		reader := kafka.NewReader(kafka.ReaderConfig{
			Brokers:     cfg.Kafka.Brokers,
			GroupID:     cfg.Kafka.GroupName,
			Topic:       dialogsmessages.TopicName,
			StartOffset: kafka.FirstOffset,
		})

		consumer, err := dialogsmessages.New(
			reader,
			do.MustInvoke[*dialogUseCases.UseCases](i),
			do.MustInvoke[*zap.Logger](i),
		)
		if err != nil {
			return nil, err
		}

		c.addShutdown(nameDialogsMessagesConsumer, sdWithoutCtx(consumer.Close))

		return consumer, nil
	})
}
