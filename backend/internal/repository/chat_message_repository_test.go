package repository

import (
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatMessageRepository(t *testing.T) {
	db := newTestDB(t)
	sessions := NewChatSessionRepository(db)
	repo := NewChatMessageRepository(db)

	sessionA := &models.ChatSession{UserID: "user_a", Title: "A", Mode: "knowledge"}
	require.NoError(t, sessions.Create(sessionA))

	sessionB := &models.ChatSession{UserID: "user_a", Title: "B", Mode: "knowledge"}
	require.NoError(t, sessions.Create(sessionB))

	t.Run("Save assigns an id", func(t *testing.T) {
		msg := &models.ChatMessage{SessionID: sessionA.ID, Role: "user", Content: "hello"}

		require.NoError(t, repo.Save(msg))

		assert.NotZero(t, msg.ID)
	})

	t.Run("GetBySessionID returns only that session's messages", func(t *testing.T) {
		require.NoError(t, repo.Save(&models.ChatMessage{SessionID: sessionA.ID, Role: "assistant", Content: "hi there"}))
		require.NoError(t, repo.Save(&models.ChatMessage{SessionID: sessionB.ID, Role: "user", Content: "other session"}))

		messagesA, err := repo.GetBySessionID(sessionA.ID)
		require.NoError(t, err)
		require.Len(t, messagesA, 2)

		contents := []string{messagesA[0].Content, messagesA[1].Content}
		assert.ElementsMatch(t, []string{"hello", "hi there"}, contents)

		for _, m := range messagesA {
			assert.Equal(t, sessionA.ID, m.SessionID)
			assert.False(t, m.CreatedAt.IsZero())
		}

		messagesB, err := repo.GetBySessionID(sessionB.ID)
		require.NoError(t, err)
		require.Len(t, messagesB, 1)
		assert.Equal(t, "other session", messagesB[0].Content)
	})

	t.Run("GetBySessionID returns an empty, non-nil slice for an unknown session", func(t *testing.T) {
		messages, err := repo.GetBySessionID(999999)

		require.NoError(t, err)
		assert.NotNil(t, messages)
		assert.Empty(t, messages)
	})
}
