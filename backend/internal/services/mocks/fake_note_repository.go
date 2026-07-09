package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeNoteRepository struct {
	Notes       []models.Note
	Note        *models.Note
	Err         error
	CreatedNote *models.Note
	UpdatedNote *models.Note
	DeletedID   int
}

func (r *FakeNoteRepository) GetAll() ([]models.Note, error) {
	return r.Notes, r.Err
}

func (r *FakeNoteRepository) GetByID(id int) (*models.Note, error) {
	return r.Note, r.Err
}

func (r *FakeNoteRepository) Create(note *models.Note) error {
	r.CreatedNote = note
	return r.Err
}

func (r *FakeNoteRepository) Update(note *models.Note) error {
	r.UpdatedNote = note
	return r.Err
}

func (r *FakeNoteRepository) Delete(id int) error {
	r.DeletedID = id
	return r.Err
}
