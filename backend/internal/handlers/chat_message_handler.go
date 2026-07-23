package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/httputil"
	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type ChatMessageService interface {
	Save(message *models.ChatMessage) error
	GetByNoteID(noteID int) ([]models.ChatMessage, error)
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

	err = h.Service.Save(&message)

	if err != nil {

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

func (h *ChatMessageHandler) GetByNoteID(
	w http.ResponseWriter,
	r *http.Request,
) {

	noteID, err := httputil.GetIDFromPath(r)

	if err != nil {

		http.Error(
			w,
			"invalid note id",
			http.StatusBadRequest,
		)

		return
	}

	messages, err := h.Service.GetByNoteID(noteID)

	if err != nil {

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
		h.GetByNoteID(w, r)

	default:
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}
