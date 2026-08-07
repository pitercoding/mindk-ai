package services

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type ChatSessionRepository interface {
	Create(session *models.ChatSession) error
	GetAll() ([]models.ChatSession, error)
	GetByID(id int) (*models.ChatSession, error)
	Update(session *models.ChatSession) error
	Delete(id int) error
}

type ChatSessionService struct {
	repo ChatSessionRepository
}

func NewChatSessionService(repo ChatSessionRepository) *ChatSessionService {
	return &ChatSessionService{
		repo: repo,
	}
}

func (s *ChatSessionService) Create(session *models.ChatSession) error {
	return s.repo.Create(session)
}

func (s *ChatSessionService) GetAll() ([]models.ChatSession, error) {
	return s.repo.GetAll()
}

func (s *ChatSessionService) GetByID(id int) (*models.ChatSession, error) {
	return s.repo.GetByID(id)
}

func (s *ChatSessionService) Update(session *models.ChatSession) error {
	return s.repo.Update(session)
}

func (s *ChatSessionService) Delete(id int) error {
	return s.repo.Delete(id)
}
