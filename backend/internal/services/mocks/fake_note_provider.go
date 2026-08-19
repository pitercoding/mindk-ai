package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeNoteProvider struct {
	Notes []models.Note
	Note  *models.Note
	Err   error

	GetAllUserID  string
	GetByIDUserID string
}

func (p *FakeNoteProvider) GetAll(userID string) ([]models.Note, error) {
	p.GetAllUserID = userID
	return p.Notes, p.Err
}

func (p *FakeNoteProvider) GetByID(id int, userID string) (*models.Note, error) {
	p.GetByIDUserID = userID
	return p.Note, p.Err
}
