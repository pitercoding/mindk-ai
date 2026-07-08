package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/services"
)

type ChatHistoryHandler struct {
	service *services.ChatHistoryService
}

func NewChatHistoryHandler(
	service *services.ChatHistoryService,
) *ChatHistoryHandler {

	return &ChatHistoryHandler{
		service: service,
	}
}

func (h *ChatHistoryHandler) GetAll(
	w http.ResponseWriter,
	r *http.Request,
) {

	history, err := h.service.GetAll()

	if err != nil {
		http.Error(
			w,
			"failed to fetch chat history",
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(history)
}
