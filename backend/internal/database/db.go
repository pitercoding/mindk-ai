package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Connect(databasePath string) error {

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

	return nil
}
