package repository

import (
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func countChunks(t *testing.T, repo *DocumentChunkRepository, documentID int) int {
	t.Helper()

	var count int
	err := repo.DB.QueryRow(
		"SELECT COUNT(*) FROM document_chunks WHERE document_id = $1",
		documentID,
	).Scan(&count)
	require.NoError(t, err)

	return count
}

func TestDocumentChunkRepository(t *testing.T) {
	db := newTestDB(t)
	documents := NewDocumentRepository(db)
	repo := NewDocumentChunkRepository(db)

	doc := &models.Document{UserID: "user_a", Name: "doc.md", Type: "md", Content: "body"}
	require.NoError(t, documents.Create(doc))

	t.Run("Create assigns an id and a creation time", func(t *testing.T) {
		chunk := &models.DocumentChunk{DocumentID: doc.ID, ChunkIndex: 0, Content: "first"}

		require.NoError(t, repo.Create(chunk))

		assert.NotZero(t, chunk.ID)
		assert.False(t, chunk.CreatedAt.IsZero())
		assert.Equal(t, 1, countChunks(t, repo, doc.ID))
	})

	t.Run("CreateMany stores every chunk with its own id", func(t *testing.T) {
		chunks := []models.DocumentChunk{
			{DocumentID: doc.ID, ChunkIndex: 1, Content: "second"},
			{DocumentID: doc.ID, ChunkIndex: 2, Content: "third"},
		}

		require.NoError(t, repo.CreateMany(chunks))

		assert.NotZero(t, chunks[0].ID)
		assert.NotZero(t, chunks[1].ID)
		assert.NotEqual(t, chunks[0].ID, chunks[1].ID)
		assert.Equal(t, 3, countChunks(t, repo, doc.ID))
	})

	t.Run("CreateMany with no chunks is a no-op", func(t *testing.T) {
		require.NoError(t, repo.CreateMany(nil))

		assert.Equal(t, 3, countChunks(t, repo, doc.ID))
	})
}
