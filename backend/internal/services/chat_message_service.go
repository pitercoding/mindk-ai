package services

import (
	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type ChatMessageRepository interface {
	Save(message *models.ChatMessage) error
	GetBySessionID(sessionID int) ([]models.ChatMessage, error)
}

// SessionOwnerChecker confirms a chat session belongs to the given user
// before its messages can be read or written. Reusing the session
// repository's own ownership check here means a session's authorization
// rule is defined in exactly one place.
type SessionOwnerChecker interface {
	GetByID(id int, userID string) (*models.ChatSession, error)
}

type ChatMessageService struct {
	repo     ChatMessageRepository
	sessions SessionOwnerChecker
}

func NewChatMessageService(
	repo ChatMessageRepository,
	sessions SessionOwnerChecker,
) *ChatMessageService {
	return &ChatMessageService{
		repo:     repo,
		sessions: sessions,
	}
}

func (s *ChatMessageService) Save(
	message *models.ChatMessage,
	userID string,
) error {

	if _, err := s.sessions.GetByID(message.SessionID, userID); err != nil {
		return err
	}

	return s.repo.Save(message)
}

func (s *ChatMessageService) GetBySessionID(
	sessionID int,
	userID string,
) ([]models.ChatMessage, error) {

	if _, err := s.sessions.GetByID(sessionID, userID); err != nil {
		return nil, err
	}

	return s.repo.GetBySessionID(sessionID)
}
