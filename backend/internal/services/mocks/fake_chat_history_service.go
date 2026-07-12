package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeChatHistoryService struct {
	SavedQuestion string
	SavedAnswer   string
	Recent        []models.ChatHistory
	LastLimit     int
	Err           error
}

func (s *FakeChatHistoryService) Save(question, answer string) error {
	s.SavedQuestion = question
	s.SavedAnswer = answer

	return s.Err
}

func (f *FakeChatHistoryService) GetRecent(limit int) ([]models.ChatHistory, error) {

	f.LastLimit = limit

	return f.Recent, f.Err
}
