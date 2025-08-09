package migrations

import (
	"context"
	"embed"
	"fmt"

	"portfolio-service/internal/config"
	"portfolio-service/internal/infrastructure/logger"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed sql/postgres/*.sql
var migrationsFS embed.FS

func RunPostgresMigrations(ctx context.Context) error {
	dsn := fmt.Sprintf(
		"pgx5://%s:%s@%s:%s/%s",
		config.AppConfig.DBUser,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBName,
	)
	l := logger.FromContext(ctx)
	files, err := migrationsFS.ReadDir("sql/postgres")
	if err != nil {
		return fmt.Errorf("failed to read embedded files: %w", err)
	}
	for _, file := range files {
		fmt.Printf("Embedded file: %s\n", file.Name())
	}
	sourceDriver, err := iofs.New(migrationsFS, "sql/postgres")
	if err != nil {
		return fmt.Errorf("failed to create source driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", sourceDriver, dsn)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	l.Infoln("Migrations applied successfully")

	return nil
}
