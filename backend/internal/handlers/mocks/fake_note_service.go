package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeNoteService struct {
	Notes []models.Note
	Note  *models.Note
	Err   error

	CreatedNote *models.Note
	UpdatedNote *models.Note
	DeletedID   int
}

func (f *FakeNoteService) Create(note *models.Note) error {
	f.CreatedNote = note
	return f.Err
}

func (f *FakeNoteService) GetAll() ([]models.Note, error) {
	return f.Notes, f.Err
}

func (f *FakeNoteService) GetByID(id int) (*models.Note, error) {
	return f.Note, f.Err
}

func (f *FakeNoteService) Update(note *models.Note) error {
	f.UpdatedNote = note
	return f.Err
}

func (f *FakeNoteService) Delete(id int) error {
	f.DeletedID = id
	return f.Err
}
