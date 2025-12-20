package internal

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/dialogs/internal/cli/deps"
)

// NewConsumeCommand команда запуска консьюмеров.
func NewConsumeCommand(container *deps.Container) *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "consume",
		Short: "run social-network dialogs consumers",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, _ []string) {
			serverLogger := container.Logger()
			cfg := container.Config()

			serverLogger.Info(fmt.Sprintf("starting consumers %s", cfg.ServiceName))

			debugServer := container.DebugServer()

			go func() {
				debugAddr := fmt.Sprintf(":%d", cfg.DebugServerListenPort)

				serverLogger.Info("starting debug server", zap.String("address", debugAddr))

				if err := debugServer.Start(debugAddr); !errors.Is(err, http.ErrServerClosed) {
					serverLogger.Error("failed to start debug server", zap.Error(err))
				}
			}()

			consumer := container.CountersMessagesConsumer()

			go consumer.Consume()

			<-cmd.Context().Done()

			serverLogger.Info("shutdown consumers")
		},
	}
}
