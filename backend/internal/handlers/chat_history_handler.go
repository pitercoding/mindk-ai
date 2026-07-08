package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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

func (h *ChatHistoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	page := 1
	limit := 10

	if value := r.URL.Query().Get("page"); value != "" {
		page, _ = strconv.Atoi(value)
	}

	if value := r.URL.Query().Get("limit"); value != "" {
		limit, _ = strconv.Atoi(value)
	}

	history, total, err := h.service.GetAll(page, limit)

	if err != nil {
		http.Error(w, "failed to fetch chat history", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"data":  history,
		"page":  page,
		"limit": limit,
		"total": total,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}
