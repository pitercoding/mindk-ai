package repository

import (
	"database/sql"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type ChatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) *ChatRepository {
	return &ChatRepository{
		db: db,
	}
}

func (r *ChatRepository) Create(history *models.ChatHistory) error {

	query := `
	INSERT INTO chat_history (
		question,
		answer
	)
	VALUES (?, ?)
	`

	result, err := r.db.Exec(
		query,
		history.Question,
		history.Answer,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	history.ID = int(id)

	return nil
}

func (r *ChatRepository) GetAll(limit, offset int) ([]models.ChatHistory, error) {

	rows, err := r.db.Query(
		`
		SELECT
			id,
			question,
			answer,
			created_at
		FROM chat_history
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
		`,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var history []models.ChatHistory

	for rows.Next() {

		var item models.ChatHistory

		err := rows.Scan(
			&item.ID,
			&item.Question,
			&item.Answer,
			&item.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		history = append(history, item)
	}

	return history, nil
}

func (r *ChatRepository) DeleteAll() error {

	_, err := r.db.Exec(`
		DELETE FROM chat_history
	`)

	return err
}

func (r *ChatRepository) Count() (int, error) {

	var total int

	err := r.db.QueryRow(
		`
		SELECT COUNT(*)
		FROM chat_history
		`,
	).Scan(&total)

	return total, err
}

func (r *ChatRepository) GetRecent(limit int) ([]models.ChatHistory, error) {

	rows, err := r.db.Query(
		`
		SELECT
			id,
			question,
			answer,
			created_at
		FROM chat_history
		ORDER BY created_at DESC
		LIMIT ?
		`,
		limit,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var history []models.ChatHistory

	for rows.Next() {

		var item models.ChatHistory

		err := rows.Scan(
			&item.ID,
			&item.Question,
			&item.Answer,
			&item.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		history = append(history, item)
	}

	return history, nil
}
