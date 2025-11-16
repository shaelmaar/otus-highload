package internal

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mongodb"
	pgMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/cobra"

	"github.com/shaelmaar/otus-highload/social-network/internal/cli/deps"
)

func NewMigrateCommand(container *deps.Container) *cobra.Command {
	var (
		mgPostgres *migrate.Migrate
		mgMongo    *migrate.Migrate
	)

	cfg := container.Config()
	pgxPool := container.PgxPool()
	logger := container.Logger()

	cmd := cobra.Command{ //nolint:exhaustruct
		Use:   "migrate",
		Short: "migrate cmd applies migrations to database",
		Long:  "migrate cmd applies migrations to database: migrate <up | down>",
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			postgresDB, err := sql.Open("pgx", pgxPool.Config().ConnString())
			if err != nil {
				return fmt.Errorf("failed to connect: %w", err)
			}

			dbInstance, _ := pgMigrate.WithInstance(postgresDB, &pgMigrate.Config{}) //nolint:exhaustruct

			mgPostgres, err = migrate.NewWithDatabaseInstance(
				"file://postgresql/migrations",
				cfg.Database.Name,
				dbInstance,
			)
			if err != nil {
				return fmt.Errorf("failed to init postgres migrations: %w", err)
			}

			mgMongo, err = migrate.New(
				"file://mongo/migrations",
				cfg.MongoDatabase.URI(),
			)
			if err != nil {
				return fmt.Errorf("failed to init mongo migrations: %w", err)
			}

			return nil
		},

		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			if mgPostgres != nil {
				_, _ = mgPostgres.Close()
			}
		},
	}

	upCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "up",
		Short: "apply up migrations",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := mgPostgres.Up()

			switch {
			case errors.Is(err, migrate.ErrNoChange):
				logger.Info("no postgres migrations found")
			case err != nil:
				return fmt.Errorf("failed to apply postgres migrations: %w", err)
			case err == nil:
				logger.Info("postgres migrations applied")
			}

			err = mgMongo.Up()
			switch {
			case errors.Is(err, migrate.ErrNoChange):
				logger.Info("no mongo migrations found")
			case err != nil:
				return fmt.Errorf("failed to apply mongo migrations: %w", err)
			case err == nil:
				logger.Info("mongo migrations applied")
			}

			return nil
		},
	}

	downCmd := &cobra.Command{ //nolint:exhaustruct
		Use:   "down",
		Short: "apply down migrations",
		RunE: func(_ *cobra.Command, _ []string) error {
			err := mgPostgres.Down()
			if err != nil {
				return fmt.Errorf("failed to down postgres migrations: %w", err)
			}

			err = mgPostgres.Down()
			if err != nil {
				return fmt.Errorf("failed to down mongo migrations: %w", err)
			}

			return nil
		},
	}

	cmd.AddCommand(upCmd)
	cmd.AddCommand(downCmd)

	return &cmd
}
