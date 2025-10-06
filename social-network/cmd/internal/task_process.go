package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/cli/deps"
)

// NewTaskProcessCommand команда запуска обработчика задач.
func NewTaskProcessCommand(container *deps.Container) *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "task_process",
		Short: "run task processing",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, _ []string) {
			serverLogger := container.Logger()
			cfg := container.Config()

			serverLogger.Info(fmt.Sprintf("starting task processing %s", cfg.ServiceName))

			errChan := make(chan error)
			defer close(errChan)

			debugServer := container.DebugServer()

			go func() {
				debugAddr := fmt.Sprintf(":%d", cfg.DebugServerListenPort)

				serverLogger.Info("starting debug server", zap.String("address", debugAddr))

				if err := debugServer.Start(debugAddr); err != nil && !errors.Is(err, http.ErrServerClosed) {
					serverLogger.Error("failed to start debug server", zap.Error(err))

					errChan <- err
				}
			}()

			userFeedTaskConsumer := container.UserFeedTaskConsumer()

			go func() {
				if err := userFeedTaskConsumer.Consume(cmd.Context()); err != nil && !errors.Is(err, context.Canceled) {
					serverLogger.Error("failed to start user feed task consumer", zap.Error(err))

					errChan <- err
				}
			}()

			userFeedChunkedTaskConsumer := container.UserFeedChunkedTaskConsumer()

			go func() {
				if err := userFeedChunkedTaskConsumer.Consume(cmd.Context()); err != nil && !errors.Is(err, context.Canceled) {
					serverLogger.Error("failed to start user feed chunked task consumer", zap.Error(err))

					errChan <- err
				}
			}()

			select {
			case <-errChan:
			case <-cmd.Context().Done():
			}

			serverLogger.Info("shutdown task processing")
		},
	}
}
