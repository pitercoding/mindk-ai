package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeChatMessageService struct {
	Messages []models.ChatMessage
	Saved    []models.ChatMessage
	Err      error
	LastSessionID int
}

func (f *FakeChatMessageService) Save(
	message *models.ChatMessage,
) error {
	f.Saved = append(
		f.Saved,
		*message,
	)

	return f.Err
}

func (f *FakeChatMessageService) GetBySessionID(
	sessionID int,
) ([]models.ChatMessage, error) {

	f.LastSessionID = sessionID
	
	return f.Messages, f.Err
}
