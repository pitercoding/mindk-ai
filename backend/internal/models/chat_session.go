package models

import "time"

type ChatSession struct {
	ID int `json:"id"`

	UserID string `json:"user_id"`

	// Session title shown in the sidebar.
	Title string `json:"title"`

	// knowledge | note
	Mode string `json:"mode"`

	// Present only for note conversations.
	NoteID *int `json:"note_id,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
