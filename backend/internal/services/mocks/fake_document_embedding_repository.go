package mocks

import (
	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type FakeDocumentEmbeddingRepository struct {
	Embeddings []models.DocumentEmbedding
	Err        error

	GetAllUserID string
}

func (f *FakeDocumentEmbeddingRepository) Create(
	embedding *models.DocumentEmbedding,
) error {

	if f.Err != nil {
		return f.Err
	}

	f.Embeddings = append(
		f.Embeddings,
		*embedding,
	)

	return nil
}

func (f *FakeDocumentEmbeddingRepository) CreateMany(
	embeddings []models.DocumentEmbedding,
) error {

	if f.Err != nil {
		return f.Err
	}

	f.Embeddings = embeddings

	return nil
}

func (f *FakeDocumentEmbeddingRepository) GetAll(userID string) (
	[]models.DocumentEmbedding,
	error,
) {
	f.GetAllUserID = userID
	return f.Embeddings, f.Err
}
