package repository

import (
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDocumentEmbeddingRepository_GetAll_UserOwnership proves that the RAG
// search path (DocumentEmbeddingRepository.GetAll, joined through
// document_chunks up to documents) never returns another user's embeddings,
// even though document_embeddings and document_chunks carry no user_id
// column of their own.
func TestDocumentEmbeddingRepository_GetAll_UserOwnership(t *testing.T) {
	db := newTestDB(t)

	documentRepo := NewDocumentRepository(db)
	chunkRepo := NewDocumentChunkRepository(db)
	embeddingRepo := NewDocumentEmbeddingRepository(db)

	const userA = "user_a"
	const userB = "user_b"

	docA := &models.Document{UserID: userA, Name: "a.txt", Type: ".txt", Content: "A content"}
	require.NoError(t, documentRepo.Create(docA))

	docB := &models.Document{UserID: userB, Name: "b.txt", Type: ".txt", Content: "B content"}
	require.NoError(t, documentRepo.Create(docB))

	chunkA := &models.DocumentChunk{DocumentID: docA.ID, ChunkIndex: 0, Content: "A chunk"}
	require.NoError(t, chunkRepo.Create(chunkA))

	chunkB := &models.DocumentChunk{DocumentID: docB.ID, ChunkIndex: 0, Content: "B chunk"}
	require.NoError(t, chunkRepo.Create(chunkB))

	require.NoError(t, embeddingRepo.Create(&models.DocumentEmbedding{
		ChunkID:   chunkA.ID,
		Embedding: `[1,0,0]`,
	}))

	require.NoError(t, embeddingRepo.Create(&models.DocumentEmbedding{
		ChunkID:   chunkB.ID,
		Embedding: `[0,1,0]`,
	}))

	embeddingsA, err := embeddingRepo.GetAll(userA)
	require.NoError(t, err)
	require.Len(t, embeddingsA, 1)
	assert.Equal(t, docA.ID, embeddingsA[0].DocumentID)
	assert.Equal(t, "A chunk", embeddingsA[0].Content)

	embeddingsB, err := embeddingRepo.GetAll(userB)
	require.NoError(t, err)
	require.Len(t, embeddingsB, 1)
	assert.Equal(t, docB.ID, embeddingsB[0].DocumentID)
	assert.Equal(t, "B chunk", embeddingsB[0].Content)
}
