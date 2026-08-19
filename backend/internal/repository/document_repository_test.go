package repository

import (
	"database/sql"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentRepository_UserOwnership(t *testing.T) {
	db := newTestDB(t)
	repo := NewDocumentRepository(db)

	const userA = "user_a"
	const userB = "user_b"

	docA := &models.Document{UserID: userA, Name: "a.txt", Type: ".txt", Content: "secret A content"}
	require.NoError(t, repo.Create(docA))

	docB := &models.Document{UserID: userB, Name: "b.txt", Type: ".txt", Content: "secret B content"}
	require.NoError(t, repo.Create(docB))

	t.Run("GetAll only returns the caller's own documents", func(t *testing.T) {
		docsA, err := repo.GetAll(userA)
		require.NoError(t, err)
		require.Len(t, docsA, 1)
		assert.Equal(t, docA.ID, docsA[0].ID)

		docsB, err := repo.GetAll(userB)
		require.NoError(t, err)
		require.Len(t, docsB, 1)
		assert.Equal(t, docB.ID, docsB[0].ID)
	})

	t.Run("GetByID denies access to another user's document", func(t *testing.T) {
		_, err := repo.GetByID(docA.ID, userB)
		assert.ErrorIs(t, err, sql.ErrNoRows)

		got, err := repo.GetByID(docA.ID, userA)
		require.NoError(t, err)
		assert.Equal(t, "a.txt", got.Name)
	})

	t.Run("Delete denies access to another user's document", func(t *testing.T) {
		err := repo.Delete(docA.ID, userB)
		assert.ErrorIs(t, err, sql.ErrNoRows)

		stillThere, err := repo.GetByID(docA.ID, userA)
		require.NoError(t, err)
		assert.NotNil(t, stillThere)
	})

	t.Run("Search only searches the caller's own documents", func(t *testing.T) {
		resultsA, err := repo.Search("secret", userA)
		require.NoError(t, err)
		require.Len(t, resultsA, 1)
		assert.Equal(t, docA.ID, resultsA[0].ID)

		resultsB, err := repo.Search("secret A", userB)
		require.NoError(t, err)
		assert.Empty(t, resultsB)
	})

	t.Run("owner can delete their own document", func(t *testing.T) {
		require.NoError(t, repo.Delete(docB.ID, userB))

		_, err := repo.GetByID(docB.ID, userB)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	})
}
