package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeDocumentRepository struct {
	Document *models.Document
	Err      error
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

func (f *FakeDocumentRepository) GetAll() ([]models.Document, error) {
	return nil, nil
}

func (f *FakeDocumentRepository) GetByID(
	id int,
) (*models.Document, error) {

	return f.Document, f.Err
}

func (f *FakeDocumentRepository) Delete(
	id int,
) error {

	return f.Err
}

func (f *FakeDocumentRepository) Search(
	query string,
) ([]models.Document, error) {

	return nil, f.Err
}
