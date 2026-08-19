package services

import (
	"database/sql"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/migrations"
	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func newChatIntegrationTestDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	// A single connection keeps every query on the same in-memory database.
	db.SetMaxOpenConns(1)

	require.NoError(t, migrations.Run(db))

	return db
}

// TestChatMessageService_CrossUserIsolation_RealDB exercises
// ChatMessageService against the real ChatSessionRepository and
// ChatMessageRepository (no fakes) to prove that chat_messages - which
// carries no user_id column of its own - is still fully isolated between
// users purely through the session ownership check enforced by real SQL.
func TestChatMessageService_CrossUserIsolation_RealDB(t *testing.T) {
	db := newChatIntegrationTestDB(t)

	sessionRepo := repository.NewChatSessionRepository(db)
	messageRepo := repository.NewChatMessageRepository(db)

	service := NewChatMessageService(messageRepo, sessionRepo)

	const userA = "user_a"
	const userB = "user_b"

	session := &models.ChatSession{UserID: userA, Title: "A's session", Mode: "knowledge"}
	require.NoError(t, sessionRepo.Create(session))

	require.NoError(t, service.Save(
		&models.ChatMessage{SessionID: session.ID, Role: "user", Content: "hi from A"},
		userA,
	))

	t.Run("owner can read their own messages", func(t *testing.T) {
		messages, err := service.GetBySessionID(session.ID, userA)
		require.NoError(t, err)
		require.Len(t, messages, 1)
		assert.Equal(t, "hi from A", messages[0].Content)
	})

	t.Run("another user cannot read messages from a session they don't own", func(t *testing.T) {
		messages, err := service.GetBySessionID(session.ID, userB)
		assert.ErrorIs(t, err, repository.ErrChatSessionNotFound)
		assert.Nil(t, messages)
	})

	t.Run("another user cannot post a message into a session they don't own", func(t *testing.T) {
		err := service.Save(
			&models.ChatMessage{SessionID: session.ID, Role: "user", Content: "injected by B"},
			userB,
		)
		assert.ErrorIs(t, err, repository.ErrChatSessionNotFound)

		messages, err := service.GetBySessionID(session.ID, userA)
		require.NoError(t, err)
		require.Len(t, messages, 1, "the injected message from B must not have been saved")
	})
}
