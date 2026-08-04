package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeDocumentService struct {
	Document  *models.Document
	Documents []models.Document
	Err       error

	Created bool
	Deleted bool
}

func (f *FakeDocumentService) Create(
	document *models.Document,
) error {

	if f.Err != nil {
		return f.Err
	}

	document.ID = 1

	f.Document = document
	f.Created = true

	return nil
}

func (f *FakeDocumentService) GetAll() ([]models.Document, error) {

	if f.Err != nil {
		return nil, f.Err
	}

	return f.Documents, nil
}

func (f *FakeDocumentService) GetByID(
	id int,
) (*models.Document, error) {

	if f.Err != nil {
		return nil, f.Err
	}

	return f.Document, nil
}

func (f *FakeDocumentService) Delete(
	id int,
) error {

	if f.Err != nil {
		return f.Err
	}

	f.Deleted = true

	return nil
}

func (f *FakeDocumentService) Search(
	query string,
) ([]models.Document, error) {

	if f.Err != nil {
		return nil, f.Err
	}

	return f.Documents, nil
}
