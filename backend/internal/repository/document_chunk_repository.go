package repository

import (
	"database/sql"
	"time"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type DocumentChunkRepository struct {
	DB *sql.DB
}

func NewDocumentChunkRepository(db *sql.DB) *DocumentChunkRepository {
	return &DocumentChunkRepository{
		DB: db,
	}
}

// Create inserts a single document chunk.
func (r *DocumentChunkRepository) Create(chunk *models.DocumentChunk) error {

	query := `
		INSERT INTO document_chunks (
			document_id,
			chunk_index,
			content
		)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.DB.QueryRow(
		query,
		chunk.DocumentID,
		chunk.ChunkIndex,
		chunk.Content,
	).Scan(&chunk.ID)

	if err != nil {
		return err
	}

	chunk.CreatedAt = time.Now()

	return nil
}

// CreateMany inserts multiple chunks.
func (r *DocumentChunkRepository) CreateMany(
	chunks []models.DocumentChunk,
) error {

	for i := range chunks {

		err := r.Create(&chunks[i])

		if err != nil {
			return err
		}
	}

	return nil
}
