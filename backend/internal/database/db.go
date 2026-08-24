package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "modernc.org/sqlite"

	"github.com/pitercoding/mindk-ai/backend/internal/config"
	"github.com/pitercoding/mindk-ai/backend/internal/migrations"
)

var (
	DB *sql.DB

	// Dialect is the migrations dialect matching the database DB is
	// connected to (migrations.DialectSQLite or migrations.DialectPostgres).
	Dialect string
)

// Connect opens the database for cfg.Environment: SQLite at cfg.DatabasePath
// in development, PostgreSQL at cfg.DatabaseURL in production.
func Connect(cfg *config.Config) error {
	if cfg.IsProduction() {
		return connectPostgres(cfg.DatabaseURL)
	}
	return connectSQLite(cfg.DatabasePath)
}

func connectSQLite(databasePath string) error {
	if dir := filepath.Dir(databasePath); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf(
				"failed to create database directory %q: %w",
				dir, err,
			)
		}
	}

	db, err := sql.Open(
		"sqlite",
		fmt.Sprintf("file:%s?_busy_timeout=5000", databasePath),
	)

	if err != nil {
		return fmt.Errorf(
			"failed to open database: %w",
			err,
		)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf(
			"failed to connect to database: %w",
			err,
		)
	}

	db.SetMaxOpenConns(1)

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return fmt.Errorf(
			"failed to enable foreign keys: %w",
			err,
		)
	}

	DB = db
	Dialect = migrations.DialectSQLite

	return nil
}

func connectPostgres(databaseURL string) error {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf(
			"failed to open database: %w",
			err,
		)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf(
			"failed to connect to database: %w",
			err,
		)
	}

	DB = db
	Dialect = migrations.DialectPostgres

	return nil
}
