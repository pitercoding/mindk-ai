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

type ChatMessageService interface {
	Save(message *models.ChatMessage, userID string) error
	GetBySessionID(sessionID int, userID string) ([]models.ChatMessage, error)
}

type ChatMessageHandler struct {
	Service ChatMessageService
}

func NewChatMessageHandler(
	service ChatMessageService,
) *ChatMessageHandler {

	return &ChatMessageHandler{
		Service: service,
	}
}

func (h *ChatMessageHandler) Save(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var message models.ChatMessage

	err := json.NewDecoder(r.Body).Decode(&message)

	if err != nil {

		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)

		return
	}

	if message.Content == "" {

		http.Error(
			w,
			"content is required",
			http.StatusBadRequest,
		)

		return
	}

	err = h.Service.Save(&message, userID)

	if err != nil {

		if errors.Is(
			err,
			repository.ErrChatSessionNotFound,
		) {

			http.Error(
				w,
				"session not found",
				http.StatusNotFound,
			)

			return
		}

		http.Error(
			w,
			"failed to save message",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(message)
}

func (h *ChatMessageHandler) GetBySessionID(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	sessionID, err := httputil.GetIDFromPath(r)

	if err != nil {

		http.Error(
			w,
			"invalid session id",
			http.StatusBadRequest,
		)

		return
	}

	messages, err := h.Service.GetBySessionID(sessionID, userID)

	if err != nil {

		if errors.Is(
			err,
			repository.ErrChatSessionNotFound,
		) {

			http.Error(
				w,
				"session not found",
				http.StatusNotFound,
			)

			return
		}

		http.Error(
			w,
			"failed to fetch chat messages",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(messages)
}

func (h *ChatMessageHandler) HandleMessages(
	w http.ResponseWriter,
	r *http.Request,
) {

	switch r.Method {

	case http.MethodPost:
		h.Save(w, r)

	case http.MethodGet:
		h.GetBySessionID(w, r)

	default:
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}
