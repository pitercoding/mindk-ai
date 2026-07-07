package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/services"
)

type ChatHandler struct {
	service *services.ChatService
}

func NewChatHandler(service *services.ChatService) *ChatHandler {
	return &ChatHandler{service: service}
}

func (h *ChatHandler) Ask(w http.ResponseWriter, r *http.Request) {
	var req models.ChatRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	answer, err := h.service.Ask(req.Message)
	if err != nil {
		log.Printf("chat request failed: %v", err)
		http.Error(w, fmt.Sprintf("failed to process chat: %v", err), http.StatusInternalServerError)
		return
	}

	resp := models.ChatResponse{
		Answer: answer,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

