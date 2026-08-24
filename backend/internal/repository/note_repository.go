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
		INSERT INTO notes (user_id, title, content, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	now := time.Now()

	err := r.DB.QueryRow(query,
		note.UserID,
		note.Title,
		note.Content,
		now,
		now,
	).Scan(&note.ID)

	if err != nil {
		return err
	}

	note.CreatedAt = now
	note.UpdatedAt = now

	return nil
}

func (r *NoteRepository) GetAll(userID string) ([]models.Note, error) {
	query := `
		SELECT
			id,
			user_id,
			title,
			content,
			created_at,
			updated_at
		FROM notes
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	notes := make([]models.Note, 0)

	for rows.Next() {
		var note models.Note

		err := rows.Scan(
			&note.ID,
			&note.UserID,
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

func (r *NoteRepository) GetByID(id int, userID string) (*models.Note, error) {
	query := `
		SELECT
			id,
			user_id,
			title,
			content,
			created_at,
			updated_at
		FROM notes
		WHERE id = $1 AND user_id = $2
	`

	var note models.Note

	err := r.DB.QueryRow(query, id, userID).Scan(
		&note.ID,
		&note.UserID,
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
		SET title = $1, content = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND user_id = $4
	`

	result, err := r.DB.Exec(query, note.Title, note.Content, note.ID, note.UserID)
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

func (r *NoteRepository) Delete(id int, userID string) error {
	query := `
		DELETE FROM notes
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
