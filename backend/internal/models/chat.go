package models

type ChatRequest struct {
	SessionID int    `json:"session_id"`
	Message   string `json:"message"`

	// Used only when session_id is omitted, to auto-create a session.
	Mode   string `json:"mode,omitempty"`
	NoteID *int   `json:"note_id,omitempty"`
	Title  string `json:"title,omitempty"`
}

type ChatSource struct {
	DocumentID int     `json:"document_id"`
	Name       string  `json:"name"`
	Score      float32 `json:"score"`
}

type ChatResponse struct {
	Answer    string       `json:"answer"`
	SessionID int          `json:"session_id"`
	Sources   []ChatSource `json:"sources,omitempty"`
}
