package models

import "time"

type ChatMessage struct {
	ID        int       `json:"id"`
	NoteID    int       `json:"note_id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}
