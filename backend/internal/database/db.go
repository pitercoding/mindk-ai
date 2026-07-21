package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Connect() error {

	db, err := sql.Open(
		"sqlite",
		"file:data/mindk.db?_busy_timeout=5000",
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

	DB = db

	return nil
}
