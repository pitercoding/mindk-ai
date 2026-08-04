package mocks

import "github.com/pitercoding/mindk-ai/backend/internal/models"

type FakeChunkCreator struct {
	Chunks []models.DocumentChunk
	Err    error
}

func (f *FakeChunkCreator) CreateMany(
	chunks []models.DocumentChunk,
) error {

	if f.Err != nil {
		return f.Err
	}

	f.Chunks = chunks

	return nil
}
