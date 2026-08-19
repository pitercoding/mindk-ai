package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeChatMessageRepository struct {
	Messages []models.ChatMessage
	Saved    *models.ChatMessage
	Err      error
}

func (f *FakeChatMessageRepository) Save(message *models.ChatMessage) error {
	f.Saved = message
	return f.Err
}

func (f *FakeChatMessageRepository) GetBySessionID(sessionID int) ([]models.ChatMessage, error) {
	return f.Messages, f.Err
}
