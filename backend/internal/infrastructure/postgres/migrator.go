// Package postgres provides PostgreSQL repository implementations.
package postgres

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies all pending up-migrations from the given directory.
// It is safe to call on every startup — already-applied migrations are skipped.
func RunMigrations(dsn, migrationsPath string) error {
	// migrationsPath must be a file:// URL, e.g. "file://migrations".
	migrator, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		dsn,
	)
	if err != nil {
		return fmt.Errorf("migrate.New: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("migrate.Up: %w", err)
	}

	return nil
}
