package internal

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/cli/deps"
	"github.com/shaelmaar/otus-highload/social-network/internal/httptransport/server"
)

// NewServeCommand команда запуска сервера.
func NewServeCommand(container *deps.Container) *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "serve",
		Short: "run social-network server",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, _ []string) {
			serverLogger := container.Logger()
			cfg := container.Config()

			serverLogger.Info(fmt.Sprintf("starting service %s", cfg.ServiceName))

			httpServer := container.HTTPServer()

			go func() {
				if err := httpServer.Serve(server.WithPort(cfg.ServerListenPort)); !errors.Is(err, http.ErrServerClosed) {
					serverLogger.Fatal("failed to start http server", zap.Error(err))
				}
			}()

			<-cmd.Context().Done()

			serverLogger.Info("shutdown service")
		},
	}
}
