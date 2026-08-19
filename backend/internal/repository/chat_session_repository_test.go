package repository

import (
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatSessionRepository_UserOwnership(t *testing.T) {
	db := newTestDB(t)
	repo := NewChatSessionRepository(db)

	const userA = "user_a"
	const userB = "user_b"

	sessionA := &models.ChatSession{UserID: userA, Title: "A's session", Mode: "knowledge"}
	require.NoError(t, repo.Create(sessionA))

	sessionB := &models.ChatSession{UserID: userB, Title: "B's session", Mode: "knowledge"}
	require.NoError(t, repo.Create(sessionB))

	t.Run("GetAll only returns the caller's own sessions", func(t *testing.T) {
		sessionsA, err := repo.GetAll(userA)
		require.NoError(t, err)
		require.Len(t, sessionsA, 1)
		assert.Equal(t, sessionA.ID, sessionsA[0].ID)

		sessionsB, err := repo.GetAll(userB)
		require.NoError(t, err)
		require.Len(t, sessionsB, 1)
		assert.Equal(t, sessionB.ID, sessionsB[0].ID)
	})

	t.Run("GetByID denies access to another user's session", func(t *testing.T) {
		_, err := repo.GetByID(sessionA.ID, userB)
		assert.ErrorIs(t, err, ErrChatSessionNotFound)

		got, err := repo.GetByID(sessionA.ID, userA)
		require.NoError(t, err)
		assert.Equal(t, "A's session", got.Title)
	})

	t.Run("Update denies access to another user's session", func(t *testing.T) {
		attempt := &models.ChatSession{ID: sessionA.ID, UserID: userB, Title: "hacked", Mode: "knowledge"}
		err := repo.Update(attempt)
		assert.ErrorIs(t, err, ErrChatSessionNotFound)

		unchanged, err := repo.GetByID(sessionA.ID, userA)
		require.NoError(t, err)
		assert.Equal(t, "A's session", unchanged.Title)
	})

	t.Run("Delete denies access to another user's session", func(t *testing.T) {
		err := repo.Delete(sessionA.ID, userB)
		assert.ErrorIs(t, err, ErrChatSessionNotFound)

		stillThere, err := repo.GetByID(sessionA.ID, userA)
		require.NoError(t, err)
		assert.NotNil(t, stillThere)
	})

	t.Run("owner can update and delete their own session", func(t *testing.T) {
		update := &models.ChatSession{ID: sessionB.ID, UserID: userB, Title: "B updated", Mode: "knowledge"}
		require.NoError(t, repo.Update(update))

		got, err := repo.GetByID(sessionB.ID, userB)
		require.NoError(t, err)
		assert.Equal(t, "B updated", got.Title)

		require.NoError(t, repo.Delete(sessionB.ID, userB))

		_, err = repo.GetByID(sessionB.ID, userB)
		assert.ErrorIs(t, err, ErrChatSessionNotFound)
	})
}
