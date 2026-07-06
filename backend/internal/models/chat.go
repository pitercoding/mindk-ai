package models

type ChatRequest struct {
	Message string `json:"message"`
}

type ChatResponse struct {
	Answer string `json:"answer"`
}
