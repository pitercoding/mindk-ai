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

func (s *ChatHistoryService) GetAll() ([]models.ChatHistory, error) {
	return s.repo.GetAll()
}

func (s *ChatHistoryService) Clear() error {
	return s.repo.DeleteAll()
}
