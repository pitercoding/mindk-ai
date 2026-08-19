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

type ChatSessionService interface {
	Create(session *models.ChatSession) error
	GetAll(userID string) ([]models.ChatSession, error)
	GetByID(id int, userID string) (*models.ChatSession, error)
	Update(session *models.ChatSession) error
	Delete(id int, userID string) error
}

type ChatSessionHandler struct {
	Service ChatSessionService
}

func NewChatSessionHandler(
	service ChatSessionService,
) *ChatSessionHandler {

	return &ChatSessionHandler{
		Service: service,
	}
}

func (h *ChatSessionHandler) HandleSessions(
	w http.ResponseWriter,
	r *http.Request,
) {

	switch r.Method {

	case http.MethodGet:
		h.GetAllSessions(w, r)

	case http.MethodPost:
		h.CreateSession(w, r)

	default:
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

func (h *ChatSessionHandler) HandleSession(
	w http.ResponseWriter,
	r *http.Request,
) {

	switch r.Method {

	case http.MethodGet:
		h.GetSessionByID(w, r)

	case http.MethodPut:
		h.UpdateSession(w, r)

	case http.MethodDelete:
		h.DeleteSession(w, r)

	default:
		http.Error(
			w,
			"method not allowed",
			http.StatusMethodNotAllowed,
		)
	}
}

func (h *ChatSessionHandler) CreateSession(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var session models.ChatSession

	err := json.NewDecoder(
		r.Body,
	).Decode(&session)

	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	if session.Title == "" || session.Mode == "" {

		http.Error(
			w,
			"title and mode are required",
			http.StatusBadRequest,
		)

		return
	}

	// The owner is always the authenticated user, never the request body
	session.UserID = userID

	err = h.Service.Create(&session)

	if err != nil {

		http.Error(
			w,
			"failed to create session",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(
		http.StatusCreated,
	)

	json.NewEncoder(w).Encode(session)
}

func (h *ChatSessionHandler) GetAllSessions(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	sessions, err := h.Service.GetAll(userID)

	if err != nil {

		http.Error(
			w,
			"failed to fetch sessions",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(sessions)
}

func (h *ChatSessionHandler) GetSessionByID(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := httputil.GetIDFromPath(r)

	if err != nil {

		http.Error(
			w,
			"invalid session id",
			http.StatusBadRequest,
		)

		return
	}

	session, err := h.Service.GetByID(id, userID)

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
			"failed to fetch session",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(session)
}

func (h *ChatSessionHandler) UpdateSession(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := httputil.GetIDFromPath(r)

	if err != nil {

		http.Error(
			w,
			"invalid session id",
			http.StatusBadRequest,
		)

		return
	}

	var session models.ChatSession

	err = json.NewDecoder(
		r.Body,
	).Decode(&session)

	if err != nil {

		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)

		return
	}

	if session.Title == "" || session.Mode == "" {

		http.Error(
			w,
			"title and mode are required",
			http.StatusBadRequest,
		)

		return
	}

	session.ID = id
	session.UserID = userID

	err = h.Service.Update(&session)

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
			"failed to update session",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(session)
}

func (h *ChatSessionHandler) DeleteSession(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := httputil.GetIDFromPath(r)

	if err != nil {

		http.Error(
			w,
			"invalid session id",
			http.StatusBadRequest,
		)

		return
	}

	err = h.Service.Delete(id, userID)

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
			"failed to delete session",
			http.StatusInternalServerError,
		)

		return
	}

	w.WriteHeader(
		http.StatusNoContent,
	)
}
