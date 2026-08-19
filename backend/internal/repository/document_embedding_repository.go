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

func (r *DocumentEmbeddingRepository) GetAll(userID string) (
	[]models.DocumentEmbedding,
	error,
) {

	query := `
		SELECT
			de.id,
			de.chunk_id,
			de.embedding,
			de.created_at,

			dc.document_id,
			dc.chunk_index,
			dc.content,

			d.name

		FROM document_embeddings de
		INNER JOIN document_chunks dc
			ON dc.id = de.chunk_id
		INNER JOIN documents d
			ON d.id = dc.document_id
		WHERE d.user_id = ?
		ORDER BY dc.document_id, dc.chunk_index
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	embeddings := make([]models.DocumentEmbedding, 0)

	for rows.Next() {

		var embedding models.DocumentEmbedding

		err := rows.Scan(
			&embedding.ID,
			&embedding.ChunkID,
			&embedding.Embedding,
			&embedding.CreatedAt,

			&embedding.DocumentID,
			&embedding.ChunkIndex,
			&embedding.Content,
			&embedding.DocumentName,
		)

		if err != nil {
			return nil, err
		}

		embeddings = append(
			embeddings,
			embedding,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return embeddings, nil
}
