package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeNoteRepository struct {
	Notes       []models.Note
	Note        *models.Note
	Err         error
	CreatedNote *models.Note
	UpdatedNote *models.Note
	DeletedID   int

	GetAllUserID  string
	GetByIDUserID string
	DeleteUserID  string
}

func (r *FakeNoteRepository) GetAll(userID string) ([]models.Note, error) {
	r.GetAllUserID = userID
	return r.Notes, r.Err
}

func (r *FakeNoteRepository) GetByID(id int, userID string) (*models.Note, error) {
	r.GetByIDUserID = userID
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

func (r *FakeNoteRepository) Delete(id int, userID string) error {
	r.DeletedID = id
	r.DeleteUserID = userID
	return r.Err
}
