package repository

import (
	"database/sql"
	"errors"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

var ErrChatSessionNotFound = errors.New("chat session not found")

type ChatSessionRepository struct {
	DB *sql.DB
}

func NewChatSessionRepository(db *sql.DB) *ChatSessionRepository {
	return &ChatSessionRepository{
		DB: db,
	}
}

func (r *ChatSessionRepository) Create(
	session *models.ChatSession,
) error {

	query := `
		INSERT INTO chat_sessions (
			user_id,
			title,
			mode,
			note_id
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	err := r.DB.QueryRow(
		query,
		session.UserID,
		session.Title,
		session.Mode,
		session.NoteID,
	).Scan(&session.ID)

	if err != nil {
		return err
	}

	created, err := r.GetByID(session.ID, session.UserID)

	if err != nil {
		return err
	}

	*session = *created

	return nil
}

func (r *ChatSessionRepository) GetAll(userID string) (
	[]models.ChatSession,
	error,
) {

	query := `
		SELECT
			id,
			user_id,
			title,
			mode,
			note_id,
			created_at,
			updated_at
		FROM chat_sessions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.DB.Query(query, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sessions := make([]models.ChatSession, 0)

	for rows.Next() {

		var session models.ChatSession

		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.Title,
			&session.Mode,
			&session.NoteID,
			&session.CreatedAt,
			&session.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		sessions = append(
			sessions,
			session,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *ChatSessionRepository) GetByID(
	id int,
	userID string,
) (*models.ChatSession, error) {

	query := `
		SELECT
			id,
			user_id,
			title,
			mode,
			note_id,
			created_at,
			updated_at
		FROM chat_sessions
		WHERE id = $1 AND user_id = $2
	`

	var session models.ChatSession

	err := r.DB.QueryRow(
		query,
		id,
		userID,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.Title,
		&session.Mode,
		&session.NoteID,
		&session.CreatedAt,
		&session.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrChatSessionNotFound
	}

	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *ChatSessionRepository) Update(
	session *models.ChatSession,
) error {

	query := `
		UPDATE chat_sessions
		SET
			title = $1,
			mode = $2,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3 AND user_id = $4
	`

	result, err := r.DB.Exec(
		query,
		session.Title,
		session.Mode,
		session.ID,
		session.UserID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrChatSessionNotFound
	}

	updated, err := r.GetByID(
		session.ID,
		session.UserID,
	)

	if err != nil {
		return err
	}

	*session = *updated

	return nil
}

func (r *ChatSessionRepository) Delete(
	id int,
	userID string,
) error {

	query := `
		DELETE FROM chat_sessions
		WHERE id = $1 AND user_id = $2
	`

	result, err := r.DB.Exec(
		query,
		id,
		userID,
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrChatSessionNotFound
	}

	return nil
}
