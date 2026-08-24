package repository

import (
	"database/sql"
	"time"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type DocumentRepository struct {
	DB *sql.DB
}

func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{DB: db}
}

func (r *DocumentRepository) Create(document *models.Document) error {
	query := `
		INSERT INTO documents (user_id, name, type, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.DB.QueryRow(
		query,
		document.UserID,
		document.Name,
		document.Type,
		document.Content,
	).Scan(&document.ID)

	if err != nil {
		return err
	}

	document.CreatedAt = time.Now()

	return nil
}

func (r *DocumentRepository) GetAll(userID string) ([]models.Document, error) {
	query := `
		SELECT
			id,
			user_id,
			name,
			type,
			content,
			created_at
		FROM documents
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documents := make([]models.Document, 0)

	for rows.Next() {
		var document models.Document

		err := rows.Scan(
			&document.ID,
			&document.UserID,
			&document.Name,
			&document.Type,
			&document.Content,
			&document.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		documents = append(documents, document)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return documents, nil
}

func (r *DocumentRepository) GetByID(id int, userID string) (*models.Document, error) {
	query := `
		SELECT
			id,
			user_id,
			name,
			type,
			content,
			created_at
		FROM documents
		WHERE id = $1 AND user_id = $2
	`

	var document models.Document

	err := r.DB.QueryRow(query, id, userID).Scan(
		&document.ID,
		&document.UserID,
		&document.Name,
		&document.Type,
		&document.Content,
		&document.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &document, nil
}

func (r *DocumentRepository) Delete(id int, userID string) error {
	query := `
	DELETE FROM documents
	WHERE id = $1 AND user_id = $2
	`

	result, err := r.DB.Exec(query, id, userID)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *DocumentRepository) Search(term string, userID string) ([]models.Document, error) {
	query := `
		SELECT
			id,
			user_id,
			name,
			type,
			content,
			created_at
		FROM documents
		WHERE user_id = $1
		AND (name LIKE $2 OR content LIKE $3)
	`

	search := "%" + term + "%"

	rows, err := r.DB.Query(
		query,
		userID,
		search,
		search,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	documents := make([]models.Document, 0)

	for rows.Next() {
		var document models.Document

		err := rows.Scan(
			&document.ID,
			&document.UserID,
			&document.Name,
			&document.Type,
			&document.Content,
			&document.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		documents = append(documents, document)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return documents, nil
}
