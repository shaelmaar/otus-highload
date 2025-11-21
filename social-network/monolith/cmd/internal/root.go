package internal

import (
	"context"
	"log"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/cli/deps"
	"github.com/shaelmaar/otus-highload/social-network/pkg/utils"
)

// Execute запуск основной команды.
func Execute(ctx context.Context) error {
	//nolint:exhaustruct // остальные поля по умолчанию
	var rootCmd = &cobra.Command{
		Use:   "social-network",
		Short: "social-network",
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	container, err := deps.New(ctx)
	utils.FatalIfErr(err)

	logger := container.Logger()

	_ = zap.ReplaceGlobals(logger)

	defer func(l *zap.Logger) {
		if err = l.Sync(); err != nil {
			log.Printf("can`t sync zap logs: %s", err)
		}
	}(logger)

	cfg := container.Config()

	//nolint:contextcheck // используется отдельный контекст для graceful shutdown.
	rootCmd.PersistentPostRun = func(_ *cobra.Command, args []string) {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		container.Shutdown(shutdownCtx)
	}

	// server.
	rootCmd.AddCommand(NewServeCommand(container))

	// websocket.
	//nolint:contextcheck // контекст берется из rootCmd.
	rootCmd.AddCommand(NewWebsocketServeCommand(container))

	// task processing.
	//nolint:contextcheck // контекст берется из rootCmd.
	rootCmd.AddCommand(NewTaskProcessCommand(container))

	// migrations.
	rootCmd.AddCommand(NewMigrateCommand(container))

	// user import command.
	//nolint:contextcheck // контекст берется из rootCmd.
	rootCmd.AddCommand(NewImportUsersCommand(container))

	//nolint:wrapcheck // не нужно оборачивать здесь ошибку.
	return rootCmd.ExecuteContext(ctx)
}
