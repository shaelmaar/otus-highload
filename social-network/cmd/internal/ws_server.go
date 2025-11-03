package internal

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/cli/deps"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/server"
)

// NewWebsocketServeCommand команда запуска вебсокет сервера.
func NewWebsocketServeCommand(container *deps.Container) *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "ws_serve",
		Short: "run social-network websocket server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, _ []string) {
			serverLogger := container.Logger()
			cfg := container.Config()

			serverLogger.Info(fmt.Sprintf("starting websocket server %s", cfg.ServiceName))

			errChan := make(chan error)
			defer close(errChan)

			debugServer := container.DebugServer()

			go func() {
				debugAddr := fmt.Sprintf(":%d", cfg.DebugServerListenPort)

				serverLogger.Info("starting debug server", zap.String("address", debugAddr))

				if err := debugServer.Start(debugAddr); !errors.Is(err, http.ErrServerClosed) {
					serverLogger.Error("failed to start debug server", zap.Error(err))

					errChan <- err
				}
			}()

			wsServer := container.WSServer()

			go func() {
				if err := wsServer.Serve(server.WithPort(cfg.WSServerListenPort)); !errors.Is(err, http.ErrServerClosed) {
					serverLogger.Error("failed to start http server", zap.Error(err))

					errChan <- err
				}
			}()

			postCreatedChunkedConsumer := container.PostCreatedChunkedTaskConsumer()

			go func() {
				if err := postCreatedChunkedConsumer.Consume(cmd.Context()); err != nil && !errors.Is(err, context.Canceled) {
					serverLogger.Error("failed to start user feed chunked task consumer", zap.Error(err))

					errChan <- err
				}
			}()

			select {
			case <-errChan:
			case <-cmd.Context().Done():
			}

			serverLogger.Info("shutdown websocket server")
		},
	}
}
