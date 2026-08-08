package services

import (
	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type ChatMessageRepository interface {
	Save(message *models.ChatMessage) error
	GetBySessionID(sessionID int) ([]models.ChatMessage, error)
}

type ChatMessageService struct {
	repo ChatMessageRepository
}

func NewChatMessageService(
	repo ChatMessageRepository,
) *ChatMessageService {
	return &ChatMessageService{
		repo: repo,
	}
}

func (s *ChatMessageService) Save(
	message *models.ChatMessage,
) error {
	return s.repo.Save(message)
}

func (s *ChatMessageService) GetBySessionID(
	sessionID int,
) ([]models.ChatMessage, error) {
	return s.repo.GetBySessionID(sessionID)
}
