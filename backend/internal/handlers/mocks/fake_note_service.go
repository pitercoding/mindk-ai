package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeNoteService struct {
	Notes []models.Note
	Note  *models.Note
	Err   error

	CreatedNote *models.Note
	UpdatedNote *models.Note
	DeletedID   int

	GetAllUserID  string
	GetByIDUserID string
	DeleteUserID  string
}

func (f *FakeNoteService) Create(note *models.Note) error {
	f.CreatedNote = note
	return f.Err
}

func (f *FakeNoteService) GetAll(userID string) ([]models.Note, error) {
	f.GetAllUserID = userID
	return f.Notes, f.Err
}

func (f *FakeNoteService) GetByID(id int, userID string) (*models.Note, error) {
	f.GetByIDUserID = userID
	return f.Note, f.Err
}

func (f *FakeNoteService) Update(note *models.Note) error {
	f.UpdatedNote = note
	return f.Err
}

func (f *FakeNoteService) Delete(id int, userID string) error {
	f.DeletedID = id
	f.DeleteUserID = userID
	return f.Err
}
