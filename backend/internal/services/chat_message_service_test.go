package services

import (
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
	"github.com/pitercoding/mindk-ai/backend/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestChatMessageServiceSave_DeniesAccessToAnotherUsersSession proves that
// saving a message into a session_id that does not belong to the caller is
// rejected, even though chat_messages carries no user_id column of its
// own - ownership flows entirely through the session it belongs to.
func TestChatMessageServiceSave_DeniesAccessToAnotherUsersSession(t *testing.T) {

	sessions := &mocks.FakeChatSessionService{
		Err: repository.ErrChatSessionNotFound,
	}

	repo := &mocks.FakeChatMessageRepository{}

	service := NewChatMessageService(repo, sessions)

	err := service.Save(
		&models.ChatMessage{SessionID: 1, Role: "user", Content: "hi"},
		"user_b",
	)

	assert.ErrorIs(t, err, repository.ErrChatSessionNotFound)
	assert.Nil(t, repo.Saved)
}

func TestChatMessageServiceSave_OwnerCanSave(t *testing.T) {

	sessions := &mocks.FakeChatSessionService{
		Session: &models.ChatSession{ID: 1, UserID: "user_a"},
	}

	repo := &mocks.FakeChatMessageRepository{}

	service := NewChatMessageService(repo, sessions)

	err := service.Save(
		&models.ChatMessage{SessionID: 1, Role: "user", Content: "hi"},
		"user_a",
	)

	require.NoError(t, err)
	require.NotNil(t, repo.Saved)
	assert.Equal(t, "hi", repo.Saved.Content)
}

func TestChatMessageServiceGetBySessionID_DeniesAccessToAnotherUsersSession(t *testing.T) {

	sessions := &mocks.FakeChatSessionService{
		Err: repository.ErrChatSessionNotFound,
	}

	repo := &mocks.FakeChatMessageRepository{
		Messages: []models.ChatMessage{
			{ID: 1, SessionID: 1, Role: "user", Content: "should not leak"},
		},
	}

	service := NewChatMessageService(repo, sessions)

	messages, err := service.GetBySessionID(1, "user_b")

	assert.ErrorIs(t, err, repository.ErrChatSessionNotFound)
	assert.Nil(t, messages)
}

func TestChatMessageServiceGetBySessionID_OwnerCanRead(t *testing.T) {

	sessions := &mocks.FakeChatSessionService{
		Session: &models.ChatSession{ID: 1, UserID: "user_a"},
	}

	repo := &mocks.FakeChatMessageRepository{
		Messages: []models.ChatMessage{
			{ID: 1, SessionID: 1, Role: "user", Content: "hi"},
		},
	}

	service := NewChatMessageService(repo, sessions)

	messages, err := service.GetBySessionID(1, "user_a")

	require.NoError(t, err)
	require.Len(t, messages, 1)
	assert.Equal(t, "hi", messages[0].Content)
}
