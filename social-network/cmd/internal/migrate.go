package internal

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate"
	pgMigrate "github.com/golang-migrate/migrate/database/postgres"
	_ "github.com/golang-migrate/migrate/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"

	"github.com/shaelmaar/otus-highload/social-network/internal/cli/deps"
)

func NewMigrateCommand(container *deps.Container) *cobra.Command {
	var mg *migrate.Migrate

	cfg := container.Config()
	pgxPool := container.PgxPool()

	cmd := cobra.Command{ //nolint:exhaustruct
		Use:   "migrate",
		Short: "migrate cmd applies migrations to database",
		Long:  "migrate cmd applies migrations to database: migrate <up | down>",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			db, err := sql.Open("pgx", pgxPool.Config().ConnString())
			if err != nil {
				return fmt.Errorf("failed to connect: %w", err)
			}

			dbInstance, _ := pgMigrate.WithInstance(db, &pgMigrate.Config{}) //nolint:exhaustruct

			mg, err = migrate.NewWithDatabaseInstance(
				"file://postgresql/migrations",
				cfg.Database.Name,
				dbInstance,
			)
			if err != nil {
				return fmt.Errorf("failed to init migrations: %w", err)
			}

			return nil
		},

		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if mg != nil {
				_, _ = mg.Close()
			}
		},
	}

	upCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "up",
		Short: "apply up migrations",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := mg.Up()

			switch {
			case errors.Is(err, migrate.ErrNoChange):
				return nil
			case err != nil:
				return fmt.Errorf("failed to apply migrations: %w", err)
			}

			return nil
		},
	}

	downCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "down",
		Short: "apply down migrations",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := mg.Down()
			if err != nil {
				return fmt.Errorf("failed to down migrations: %w", err)
			}

			return nil
		},
	}

	cmd.AddCommand(upCmd)
	cmd.AddCommand(downCmd)

	return &cmd
}
