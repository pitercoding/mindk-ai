package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeDocumentRepository struct {
	Document *models.Document
	Err      error

	GetAllUserID  string
	GetByIDUserID string
	DeleteUserID  string
	SearchUserID  string
}

func (f *FakeDocumentRepository) Create(
	document *models.Document,
) error {

	if f.Err != nil {
		return f.Err
	}

	document.ID = 1

	f.Document = document

	return nil
}

func (f *FakeDocumentRepository) GetAll(userID string) ([]models.Document, error) {
	f.GetAllUserID = userID
	return nil, nil
}

func (f *FakeDocumentRepository) GetByID(
	id int,
	userID string,
) (*models.Document, error) {

	f.GetByIDUserID = userID

	return f.Document, f.Err
}

func (f *FakeDocumentRepository) Delete(
	id int,
	userID string,
) error {

	f.DeleteUserID = userID

	return f.Err
}

func (f *FakeDocumentRepository) Search(
	query string,
	userID string,
) ([]models.Document, error) {

	f.SearchUserID = userID

	return nil, f.Err
}
