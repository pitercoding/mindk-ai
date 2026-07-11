package services

import (
	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type ChatHistoryRepository interface {
	Create(history *models.ChatHistory) error
	GetAll(limit, offset int) ([]models.ChatHistory, error)
	GetRecent(limit int) ([]models.ChatHistory, error)
	Count() (int, error)
	DeleteAll() error
}

type ChatHistoryService struct {
	repo ChatHistoryRepository
}

func NewChatHistoryService(repo ChatHistoryRepository) *ChatHistoryService {

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

func (s *ChatHistoryService) GetRecent(limit int) ([]models.ChatHistory, error) {
	return s.repo.GetRecent(limit)
}

func (s *ChatHistoryService) Clear() error {
	return s.repo.DeleteAll()
}
