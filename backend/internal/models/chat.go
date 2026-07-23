package models

type ChatRequest struct {
	Message string       `json:"message"`
	Context *ChatContext `json:"context,omitempty"`
}

type ChatContext struct {
	NoteID  int    `json:"note_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Answer string `json:"answer"`
}
