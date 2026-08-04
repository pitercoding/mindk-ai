package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeEmbeddingService struct {
	Called bool
	Err    error
}

func (f *FakeEmbeddingService) GenerateForChunks(
	chunks []models.DocumentChunk,
) error {

	f.Called = true

	return f.Err
}
