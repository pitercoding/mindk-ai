package services

import (
	"errors"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/services/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentEmbeddingServiceGenerateForChunks(
	t *testing.T,
) {

	t.Run(
		"generates and saves embeddings successfully",
		func(t *testing.T) {

			repository := &mocks.FakeDocumentEmbeddingRepository{}

			generator := &mocks.FakeEmbeddingGenerator{
				Vector: []float32{
					0.1,
					0.2,
					0.3,
				},
			}

			service := NewDocumentEmbeddingService(
				repository,
				generator,
			)

			chunks := []models.DocumentChunk{
				{
					ID:      1,
					Content: "Go is a programming language",
				},
				{
					ID:      2,
					Content: "AI uses embeddings",
				},
			}

			err := service.GenerateForChunks(chunks)

			require.NoError(
				t,
				err,
			)

			assert.True(
				t,
				generator.Called,
			)

			assert.Len(
				t,
				repository.Embeddings,
				2,
			)

			assert.Equal(
				t,
				1,
				repository.Embeddings[0].ChunkID,
			)

			assert.NotEmpty(
				t,
				repository.Embeddings[0].Embedding,
			)
		},
	)

	t.Run(
		"returns error when embedding generation fails",
		func(t *testing.T) {

			repository := &mocks.FakeDocumentEmbeddingRepository{}

			generator := &mocks.FakeEmbeddingGenerator{
				Err: errors.New("embedding failed"),
			}

			service := NewDocumentEmbeddingService(
				repository,
				generator,
			)

			chunks := []models.DocumentChunk{
				{
					ID:      1,
					Content: "test",
				},
			}

			err := service.GenerateForChunks(chunks)

			assert.Error(
				t,
				err,
			)

			assert.Empty(
				t,
				repository.Embeddings,
			)
		},
	)
}
