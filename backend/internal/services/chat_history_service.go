package services

import (
	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
)

type ChatHistoryService struct {
	repo *repository.ChatRepository
}

func NewChatHistoryService(repo *repository.ChatRepository) *ChatHistoryService {

	return &ChatHistoryService{
		repo: repo,
	}
}

func (s *ChatHistoryService) Save(question, answer string) error {

	history := &models.ChatHistory{
		Question: question,
		Answer:   answer,
	}

	return s.repo.Create(history)
}

func (s *ChatHistoryService) GetAll(page, limit int) ([]models.ChatHistory, int, error) {

	offset := (page - 1) * limit

	history, err := s.repo.GetAll(limit, offset)

	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count()

	if err != nil {
		return nil, 0, err
	}

	return history, total, nil
}

func (s *ChatHistoryService) Clear() error {
	return s.repo.DeleteAll()
}
