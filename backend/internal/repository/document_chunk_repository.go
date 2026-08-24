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

// GetByDocumentID returns all chunks that belong to a document.
func (r *DocumentChunkRepository) GetByDocumentID(
	documentID int,
) ([]models.DocumentChunk, error) {

	query := `
		SELECT
			id,
			document_id,
			chunk_index,
			content,
			created_at
		FROM document_chunks
		WHERE document_id = $1
		ORDER BY chunk_index
	`

	rows, err := r.DB.Query(
		query,
		documentID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	chunks := make([]models.DocumentChunk, 0)

	for rows.Next() {

		var chunk models.DocumentChunk

		err := rows.Scan(
			&chunk.ID,
			&chunk.DocumentID,
			&chunk.ChunkIndex,
			&chunk.Content,
			&chunk.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		chunks = append(
			chunks,
			chunk,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return chunks, nil
}

// DeleteByDocumentID removes all chunks of a document.
func (r *DocumentChunkRepository) DeleteByDocumentID(
	documentID int,
) error {

	query := `
		DELETE FROM document_chunks
		WHERE document_id = $1
	`

	_, err := r.DB.Exec(
		query,
		documentID,
	)

	return err
}
