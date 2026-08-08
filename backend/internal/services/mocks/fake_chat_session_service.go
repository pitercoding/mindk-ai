package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeChatSessionService struct {
	Session *models.ChatSession
	Err     error
	Created []models.ChatSession
}

func (f *FakeChatSessionService) Create(
	session *models.ChatSession,
) error {

	f.Created = append(
		f.Created,
		*session,
	)

	if f.Err != nil {
		return f.Err
	}

	if f.Session != nil {
		session.ID = f.Session.ID
	}

	return nil
}

func (f *FakeChatSessionService) GetByID(
	id int,
) (*models.ChatSession, error) {

	if f.Err != nil {
		return nil, f.Err
	}

	return f.Session, nil
}
