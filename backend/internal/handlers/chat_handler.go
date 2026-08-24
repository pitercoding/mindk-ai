package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/auth"
	"github.com/pitercoding/mindk-ai/backend/internal/httputil"
	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
)

const maxChatMessageLength = 8_000

type ChatService interface {
	Ask(
		userID string,
		sessionID int,
		message string,
		mode string,
		noteID *int,
		title string,
	) (string, int, []models.ChatSource, error)
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
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.ChatRequest

	if err := httputil.DecodeJSON(w, r, httputil.MaxJSONBodyBytes, &req); err != nil {
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

	if len(req.Message) > maxChatMessageLength {
		http.Error(
			w,
			"message exceeds maximum length",
			http.StatusBadRequest,
		)
		return
	}

	// Mode is only used to auto-create a session, so it only needs
	// validating on that path - an existing session already has a mode.
	if req.SessionID == 0 && req.Mode != "knowledge" && req.Mode != "note" {
		http.Error(
			w,
			"mode must be 'knowledge' or 'note'",
			http.StatusBadRequest,
		)
		return
	}

	answer, sessionID, sources, err := h.service.Ask(
		userID,
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

		httputil.LogError(r, "chat request failed", err)

		http.Error(
			w,
			"failed to process chat request",
			http.StatusInternalServerError,
		)

		return
	}

	resp := models.ChatResponse{
		Answer:    answer,
		SessionID: sessionID,
		Sources:   sources,
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(resp)

}
