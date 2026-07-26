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
		INSERT INTO documents (name, type, content, created_at)
		VALUES (?, ?, ?, ?)
	`
	now := time.Now()

	result, err := r.DB.Exec(
		query,
		document.Name,
		document.Type,
		document.Content,
		now,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	document.ID = int(id)
	document.CreatedAt = now

	return nil
}

func (r *DocumentRepository) GetAll() ([]models.Document, error) {
	query := `
		SELECT
			id,
			name,
			type,
			content,
			created_at
		FROM documents
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []models.Document

	for rows.Next() {
		var document models.Document

		err := rows.Scan(
			&document.ID,
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

func (r *DocumentRepository) GetByID(id int) (*models.Document, error) {
	query := `
		SELECT
			id,
			name,
			type,
			content,
			created_at
		FROM documents
		WHERE id = ?
	`

	var document models.Document

	err := r.DB.QueryRow(query, id).Scan(
		&document.ID,
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

func (r *DocumentRepository) Delete(id int) error {
	query := `
	DELETE FROM documents
	WHERE id = ?
	`

	_, err := r.DB.Exec(query, id)
	return err
}

func (r *DocumentRepository) Search(term string) ([]models.Document, error) {
	query := `
		SELECT
			id,
			name,
			type,
			content,
			created_at
		FROM documents
		WHERE name LIKE ?
		OR content LIKE ?
	`

	search := "%" + term + "%"

	rows, err := r.DB.Query(
		query,
		search,
		search,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []models.Document

	for rows.Next() {
		var document models.Document

		err := rows.Scan(
			&document.ID,
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
