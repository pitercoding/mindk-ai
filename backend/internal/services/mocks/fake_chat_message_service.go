package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeChatMessageService struct {
	Messages      []models.ChatMessage
	Saved         []models.ChatMessage
	Err           error
	LastSessionID int

	LastUserID string
}

func (f *FakeChatMessageService) Save(
	message *models.ChatMessage,
	userID string,
) error {
	f.Saved = append(
		f.Saved,
		*message,
	)

	f.LastUserID = userID

	return f.Err
}

func (f *FakeChatMessageService) GetBySessionID(
	sessionID int,
	userID string,
) ([]models.ChatMessage, error) {

	f.LastSessionID = sessionID
	f.LastUserID = userID

	return f.Messages, f.Err
}
