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

func (r *ChatMessageRepository) Save(message *models.ChatMessage) error {

	query := `
		INSERT INTO chat_messages (
			note_id,
			role,
			content
		)
		VALUES (?, ?, ?)
	`

	result, err := r.DB.Exec(
		query,
		message.NoteID,
		message.Role,
		message.Content,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	message.ID = int(id)

	return nil
}

func (r *ChatMessageRepository) GetByNoteID(
	noteID int,
) ([]models.ChatMessage, error) {

	query := `
		SELECT
			id,
			note_id,
			role,
			content,
			created_at
		FROM chat_messages
		WHERE note_id = ?
		ORDER BY created_at ASC
	`

	rows, err := r.DB.Query(query, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.ChatMessage

	for rows.Next() {

		var message models.ChatMessage

		err := rows.Scan(
			&message.ID,
			&message.NoteID,
			&message.Role,
			&message.Content,
			&message.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return messages, nil
}

