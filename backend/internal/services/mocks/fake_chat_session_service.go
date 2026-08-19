package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeChatSessionService struct {
	Session  *models.ChatSession
	Sessions []models.ChatSession
	Created  []models.ChatSession
	Updated  []models.ChatSession
	Deleted  []int
	Err      error

	GetAllUserID  string
	GetByIDUserID string
	DeleteUserID  string
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

func (f *FakeChatSessionService) GetAll(userID string) (
	[]models.ChatSession,
	error,
) {

	f.GetAllUserID = userID

	if f.Err != nil {
		return nil, f.Err
	}

	return f.Sessions, nil
}

func (f *FakeChatSessionService) GetByID(
	id int,
	userID string,
) (*models.ChatSession, error) {

	f.GetByIDUserID = userID

	if f.Err != nil {
		return nil, f.Err
	}

	return f.Session, nil
}

func (f *FakeChatSessionService) Update(
	session *models.ChatSession,
) error {

	f.Updated = append(
		f.Updated,
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

func (f *FakeChatSessionService) Delete(
	id int,
	userID string,
) error {

	f.Deleted = append(
		f.Deleted,
		id,
	)

	f.DeleteUserID = userID

	return f.Err
}
