package repository

import (
	"database/sql"
	"time"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type NoteRepository struct {
	DB *sql.DB
}

func NewNoteRepository(db *sql.DB) *NoteRepository {
	return &NoteRepository{DB: db}
}

func (r *NoteRepository) Create(note *models.Note) error {
	query := `
		INSERT INTO notes (title, content, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`

	now := time.Now()

	result, err := r.DB.Exec(query,
		note.Title,
		note.Content,
		now,
		now,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	note.ID = int(id)
	note.CreatedAt = now
	note.UpdatedAt = now

	return nil
}

func (r *NoteRepository) GetAll() ([]models.Note, error) {
	query := `
		SELECT
			id,
			title,
			content,
			created_at,
			updated_at
		FROM notes
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notes []models.Note

	for rows.Next() {
		var note models.Note

		err := rows.Scan(
			&note.ID,
			&note.Title,
			&note.Content,
			&note.CreatedAt,
			&note.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		notes = append(notes, note)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notes, nil
}

func (r *NoteRepository) GetByID(id int) (*models.Note, error) {
	query := `
		SELECT
			id,
			title,
			content,
			created_at,
			updated_at
		FROM notes
		WHERE id = ?
	`

	var note models.Note

	err := r.DB.QueryRow(query, id).Scan(
		&note.ID,
		&note.Title,
		&note.Content,
		&note.CreatedAt,
		&note.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &note, nil
}

func (r *NoteRepository) Update(note *models.Note) error {
	query := `
		UPDATE notes
		SET title = ?, content = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := r.DB.Exec(query, note.Title, note.Content, note.ID)
	return err
}

func (r *NoteRepository) Delete(id int) error {
	query := `
		DELETE FROM notes
		WHERE id = ?
	`

	_, err := r.DB.Exec(query, id)
	return err
}

