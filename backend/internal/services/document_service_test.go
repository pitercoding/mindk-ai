package services

import (
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeEmbeddingCreator struct {
	Called bool
	Chunks []models.DocumentChunk
	Err    error
}

func (f *fakeEmbeddingCreator) GenerateForChunks(
	chunks []models.DocumentChunk,
) error {

	f.Called = true
	f.Chunks = chunks

	return f.Err
}

func TestDocumentServiceCreate(t *testing.T) {

	documentRepo := &mocks.FakeDocumentRepository{}

	chunkCreator := &mocks.FakeChunkCreator{}

	embeddingCreator := &fakeEmbeddingCreator{}

	service := NewDocumentService(
		documentRepo,
		chunkCreator,
		embeddingCreator,
	)

	document := &models.Document{
		Name: "test.txt",
		Content: `
			This is a test document.
			It contains enough content
			to generate chunks.
		`,
	}

	err := service.Create(document)

	require.NoError(t, err)

	// document was created
	assert.Equal(
		t,
		1,
		document.ID,
	)

	// chunks were created
	assert.NotEmpty(
		t,
		chunkCreator.Chunks,
	)

	assert.Equal(
		t,
		1,
		chunkCreator.Chunks[0].DocumentID,
	)

	// embeddings were generated
	assert.True(
		t,
		embeddingCreator.Called,
	)

	assert.Equal(
		t,
		len(chunkCreator.Chunks),
		len(embeddingCreator.Chunks),
	)
}
