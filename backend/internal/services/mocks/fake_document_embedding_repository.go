package mocks

import (
	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type FakeDocumentEmbeddingRepository struct {
	Embeddings []models.DocumentEmbedding
	Err        error
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

func (f *FakeDocumentEmbeddingRepository) GetByChunkID(
	chunkID int,
) (*models.DocumentEmbedding, error) {

	if f.Err != nil {
		return nil, f.Err
	}

	for _, embedding := range f.Embeddings {

		if embedding.ChunkID == chunkID {
			return &embedding, nil
		}
	}

	return nil, nil
}

func (f *FakeDocumentEmbeddingRepository) DeleteByChunkID(
	chunkID int,
) error {

	return f.Err
}

func (f *FakeDocumentEmbeddingRepository) GetAll() (
	[]models.DocumentEmbedding,
	error,
) {
	return f.Embeddings, f.Err
}
