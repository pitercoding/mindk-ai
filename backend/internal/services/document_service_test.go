package services

import (
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentServiceCreate(t *testing.T) {

	documentRepo := &mocks.FakeDocumentRepository{}

	chunkService := &mocks.FakeChunkCreator{}

	embeddingService := &mocks.FakeEmbeddingService{}

	service := NewDocumentService(
		documentRepo,
		chunkService,
		embeddingService,
	)

	document := &models.Document{
		Name: "test.txt",
		Type: "text/plain",
		Content: `
			This is a test document.
			It contains enough content
			to generate chunks.
		`,
	}

	err := service.Create(document)

	require.NoError(t, err)

	t.Run("document was created", func(t *testing.T) {

		assert.Equal(
			t,
			1,
			document.ID,
		)

		assert.Equal(
			t,
			document,
			documentRepo.Document,
		)
	})

	t.Run("chunks were created", func(t *testing.T) {

		assert.NotEmpty(
			t,
			chunkService.Chunks,
		)

		assert.Equal(
			t,
			1,
			chunkService.Chunks[0].DocumentID,
		)
	})

	t.Run("embeddings were generated", func(t *testing.T) {

		assert.True(
			t,
			embeddingService.Called,
		)
	})
}
