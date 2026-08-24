package repository

import (
	"database/sql"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type ChatMessageRepository struct {
	DB *sql.DB
}

func NewChatMessageRepository(db *sql.DB) *ChatMessageRepository {
	return &ChatMessageRepository{
		DB: db,
	}
}

func (r *ChatMessageRepository) Save(
	message *models.ChatMessage,
) error {

	query := `
		INSERT INTO chat_messages (
			session_id,
			role,
			content
		)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	err := r.DB.QueryRow(
		query,
		message.SessionID,
		message.Role,
		message.Content,
	).Scan(&message.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *ChatMessageRepository) GetBySessionID(
	sessionID int,
) ([]models.ChatMessage, error) {

	query := `
		SELECT
			id,
			session_id,
			role,
			content,
			created_at
		FROM chat_messages
		WHERE session_id = $1
		ORDER BY created_at ASC
	`

	rows, err := r.DB.Query(
		query,
		sessionID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	messages := make([]models.ChatMessage, 0)

	for rows.Next() {

		var message models.ChatMessage

		err := rows.Scan(
			&message.ID,
			&message.SessionID,
			&message.Role,
			&message.Content,
			&message.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		messages = append(
			messages,
			message,
		)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}
