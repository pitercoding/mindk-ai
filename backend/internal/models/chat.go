package models

type ChatRequest struct {
	SessionID int    `json:"session_id"`
	Message   string `json:"message"`
}

type ChatResponse struct {
	Answer string `json:"answer"`
}
