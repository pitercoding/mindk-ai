package repository

import (
	"database/sql"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/migrations"
	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	// A single connection keeps every query on the same in-memory database.
	db.SetMaxOpenConns(1)

	require.NoError(t, migrations.Run(db))

	return db
}

func TestNoteRepository_UserOwnership(t *testing.T) {
	db := newTestDB(t)
	repo := NewNoteRepository(db)

	const userA = "user_a"
	const userB = "user_b"

	noteA := &models.Note{UserID: userA, Title: "A's note", Content: "secret A"}
	require.NoError(t, repo.Create(noteA))

	noteB := &models.Note{UserID: userB, Title: "B's note", Content: "secret B"}
	require.NoError(t, repo.Create(noteB))

	t.Run("GetAll only returns the caller's own notes", func(t *testing.T) {
		notesA, err := repo.GetAll(userA)
		require.NoError(t, err)
		require.Len(t, notesA, 1)
		assert.Equal(t, noteA.ID, notesA[0].ID)

		notesB, err := repo.GetAll(userB)
		require.NoError(t, err)
		require.Len(t, notesB, 1)
		assert.Equal(t, noteB.ID, notesB[0].ID)
	})

	t.Run("GetByID denies access to another user's note", func(t *testing.T) {
		_, err := repo.GetByID(noteA.ID, userB)
		assert.ErrorIs(t, err, sql.ErrNoRows)

		got, err := repo.GetByID(noteA.ID, userA)
		require.NoError(t, err)
		assert.Equal(t, "A's note", got.Title)
	})

	t.Run("Update denies access to another user's note", func(t *testing.T) {
		attempt := &models.Note{ID: noteA.ID, UserID: userB, Title: "hacked", Content: "hacked"}
		err := repo.Update(attempt)
		assert.ErrorIs(t, err, sql.ErrNoRows)

		unchanged, err := repo.GetByID(noteA.ID, userA)
		require.NoError(t, err)
		assert.Equal(t, "A's note", unchanged.Title)
	})

	t.Run("Delete denies access to another user's note", func(t *testing.T) {
		err := repo.Delete(noteA.ID, userB)
		assert.ErrorIs(t, err, sql.ErrNoRows)

		stillThere, err := repo.GetByID(noteA.ID, userA)
		require.NoError(t, err)
		assert.NotNil(t, stillThere)
	})

	t.Run("owner can update and delete their own note", func(t *testing.T) {
		update := &models.Note{ID: noteB.ID, UserID: userB, Title: "B updated", Content: "still secret"}
		require.NoError(t, repo.Update(update))

		got, err := repo.GetByID(noteB.ID, userB)
		require.NoError(t, err)
		assert.Equal(t, "B updated", got.Title)

		require.NoError(t, repo.Delete(noteB.ID, userB))

		_, err = repo.GetByID(noteB.ID, userB)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}
