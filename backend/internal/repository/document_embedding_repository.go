package repository

import (
	"database/sql"
	"time"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type DocumentEmbeddingRepository struct {
	DB *sql.DB
}

func NewDocumentEmbeddingRepository(db *sql.DB) *DocumentEmbeddingRepository {
	return &DocumentEmbeddingRepository{
		DB: db,
	}
}

func (r *DocumentEmbeddingRepository) Create(embedding *models.DocumentEmbedding,
) error {

	query := `
		INSERT INTO document_embeddings (
			chunk_id,
			embedding
		)
		VALUES (?, ?)
	`

	result, err := r.DB.Exec(
		query,
		embedding.ChunkID,
		embedding.Embedding,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	embedding.ID = int(id)
	embedding.CreatedAt = time.Now()

	return nil
}

func (r *DocumentEmbeddingRepository) CreateMany(
	embeddings []models.DocumentEmbedding,
) error {

	for i := range embeddings {

		err := r.Create(&embeddings[i])

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *DocumentEmbeddingRepository) GetByChunkID(
	chunkID int,
) (*models.DocumentEmbedding, error) {

	query := `
        SELECT
            id,
            chunk_id,
            embedding,
            created_at
        FROM document_embeddings
        WHERE chunk_id = ?
    `

	var embedding models.DocumentEmbedding

	err := r.DB.QueryRow(query, chunkID).Scan(
		&embedding.ID,
		&embedding.ChunkID,
		&embedding.Embedding,
		&embedding.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &embedding, nil
}

func (r *DocumentEmbeddingRepository) DeleteByChunkID(
	chunkID int,
) error {

	query := `
		DELETE FROM document_embeddings
		WHERE chunk_id = ?
	`

	_, err := r.DB.Exec(
		query,
		chunkID,
	)

	return err
}
