package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeNoteProvider struct {
	Notes []models.Note
	Err   error
}

func (p *FakeNoteProvider) GetAll() ([]models.Note, error) {
	return p.Notes, p.Err
}
