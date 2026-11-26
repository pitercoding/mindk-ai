package migrations

import (
	"database/sql"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func openSQLite(t *testing.T) *sql.DB {
	t.Helper()

	path := filepath.Join(t.TempDir(), "migrations-test.db")

	db, err := sql.Open("sqlite", "file:"+path)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	db.SetMaxOpenConns(1)

	return db
}

func tableExists(t *testing.T, db *sql.DB, name string) bool {
	t.Helper()

	var found string
	err := db.QueryRow(
		"SELECT name FROM sqlite_master WHERE type = 'table' AND name = ?",
		name,
	).Scan(&found)

	if err == sql.ErrNoRows {
		return false
	}
	require.NoError(t, err)

	return true
}

func TestRun_SQLite(t *testing.T) {
	t.Run("creates every application table", func(t *testing.T) {
		db := openSQLite(t)

		require.NoError(t, Run(db, DialectSQLite))

		for _, table := range []string{
			"notes",
			"documents",
			"document_chunks",
			"document_embeddings",
			"chat_sessions",
			"chat_messages",
		} {
			assert.True(t, tableExists(t, db, table), "table %q should exist", table)
		}
	})

	t.Run("is idempotent when run again", func(t *testing.T) {
		db := openSQLite(t)

		require.NoError(t, Run(db, DialectSQLite))
		require.NoError(t, Run(db, DialectSQLite))
	})

	t.Run("adds user_id to the owned tables", func(t *testing.T) {
		db := openSQLite(t)
		require.NoError(t, Run(db, DialectSQLite))

		for _, table := range []string{"notes", "documents", "chat_sessions"} {
			_, err := db.Exec("SELECT user_id FROM " + table)
			assert.NoError(t, err, "%s should have a user_id column", table)
		}
	})
}

func TestRun_UnknownDialect(t *testing.T) {
	db := openSQLite(t)

	err := Run(db, "mysql")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown migrations dialect")
}
