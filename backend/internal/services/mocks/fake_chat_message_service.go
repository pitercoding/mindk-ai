package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeChatMessageService struct {
	Messages []models.ChatMessage
	Saved    []models.ChatMessage
	Err      error
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

func (f *FakeChatMessageService) GetByNoteID(
	noteID int,
) ([]models.ChatMessage, error) {

	return f.Messages, f.Err
}
