package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
)

type ChatService interface {
	Ask(
		sessionID int,
		message string,
		mode string,
		noteID *int,
		title string,
	) (string, int, error)
}

type ChatHandler struct {
	service ChatService
}

func NewChatHandler(service ChatService) *ChatHandler {
	return &ChatHandler{
		service: service,
	}
}

func (h *ChatHandler) Ask(w http.ResponseWriter, r *http.Request) {
	var req models.ChatRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(
			w,
			"invalid request",
			http.StatusBadRequest,
		)
		return
	}

	if req.Message == "" {
		http.Error(
			w,
			"message is required",
			http.StatusBadRequest,
		)
		return
	}

	answer, sessionID, err := h.service.Ask(
		req.SessionID,
		req.Message,
		req.Mode,
		req.NoteID,
		req.Title,
	)

	if err != nil {

		if errors.Is(err, repository.ErrChatSessionNotFound) {

			http.Error(
				w,
				"session not found",
				http.StatusNotFound,
			)

			return
		}

		log.Printf("chat request failed: %v", err)

		http.Error(
			w,
			fmt.Sprintf("failed to process chat: %v", err),
			http.StatusInternalServerError,
		)

		return
	}

	resp := models.ChatResponse{
		Answer:    answer,
		SessionID: sessionID,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(resp)

}
