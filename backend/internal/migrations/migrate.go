package migrations

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	migratepgx "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

const (
	DialectSQLite   = "sqlite"
	DialectPostgres = "postgres"
)

//go:embed sqlite/*.sql
var sqliteFS embed.FS

//go:embed postgres/*.sql
var postgresFS embed.FS

// Run applies all pending migrations to db. dialect selects both the SQL
// dialect of the migration files and the golang-migrate database driver used
// to apply them - it must be DialectSQLite or DialectPostgres.
func Run(db *sql.DB, dialect string) error {
	var (
		migrationsFS embed.FS
		subdir       string
		driver       database.Driver
		err          error
	)

	switch dialect {
	case DialectSQLite:
		migrationsFS, subdir = sqliteFS, "sqlite"
		driver, err = sqlite.WithInstance(db, &sqlite.Config{})
	case DialectPostgres:
		migrationsFS, subdir = postgresFS, "postgres"
		driver, err = migratepgx.WithInstance(db, &migratepgx.Config{})
	default:
		return fmt.Errorf("unknown migrations dialect %q", dialect)
	}
	if err != nil {
		return fmt.Errorf("failed to create %s driver: %w", dialect, err)
	}

	src, err := fs.Sub(migrationsFS, subdir)
	if err != nil {
		return fmt.Errorf("failed to select %s migrations: %w", dialect, err)
	}

	d, err := iofs.New(src, ".")
	if err != nil {
		return fmt.Errorf("failed to create iofs driver: %w", err)
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		d,
		dialect,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}
