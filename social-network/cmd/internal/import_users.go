package internal

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"github.com/shaelmaar/otus-highload/social-network/internal/cli/deps"
)

// NewImportUsersCommand команда запуска сервера.
func NewImportUsersCommand(container *deps.Container) *cobra.Command {
	return &cobra.Command{ //nolint:exhaustruct
		Use:   "import_users [filename]",
		Short: "import users from csv file",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			filename := args[0]

			cmdLogger := container.Logger()
			useCases := container.UserUseCases()

			cmdLogger.Info("start importing users")

			now := time.Now()

			defer func() {
				cmdLogger.Info("finished importing users", zap.Duration("elapsed", time.Since(now)))
			}()

			err := useCases.ImportUsers(cmd.Context(), filename)
			if err != nil {
				return fmt.Errorf("failed to import users: %w", err)
			}

			return nil
		},
	}
}
