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

const maxChatSessionTitleLength = 200

// validateChatSession returns a client-facing error message, or "" if
// session is valid. Shared by CreateSession and UpdateSession, which apply
// the exact same rules to the client-supplied fields.
func validateChatSession(session models.ChatSession) string {
	if session.Title == "" || session.Mode == "" {
		return "title and mode are required"
	}

	if len(session.Title) > maxChatSessionTitleLength {
		return "title exceeds maximum length"
	}

	if session.Mode != "knowledge" && session.Mode != "note" {
		return "mode must be 'knowledge' or 'note'"
	}

	return ""
}

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

// CreateSession creates a chat session owned by the authenticated user.
//
//	@Summary		Create chat session
//	@Description	Creates a chat session. mode must be "knowledge" or "note"; note_id is meaningful only for "note" mode. user_id is taken from the token, not the body.
//	@Tags			chat-sessions
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			session	body		models.ChatSession	true	"Session to create (title and mode are required)"
//	@Success		201		{object}	models.ChatSession
//	@Failure		400		{string}	string	"title and mode are required, title too long, or mode is invalid"
//	@Failure		401		{string}	string	"unauthorized"
//	@Failure		413		{string}	string	"request body too large"
//	@Failure		500		{string}	string	"failed to create session"
//	@Router			/chat/sessions [post]
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

	if err := httputil.DecodeJSON(w, r, httputil.MaxJSONBodyBytes, &session); err != nil {
		return
	}

	if err := validateChatSession(session); err != "" {
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	// The owner is always the authenticated user, never the request body
	session.UserID = userID

	err := h.Service.Create(&session)

	if err != nil {

		httputil.LogError(r, "failed to create chat session", err)

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

// GetAllSessions lists every chat session owned by the authenticated user.
//
//	@Summary		List chat sessions
//	@Description	Returns every chat session owned by the authenticated user.
//	@Tags			chat-sessions
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.ChatSession
//	@Failure		401	{string}	string	"unauthorized"
//	@Failure		500	{string}	string	"failed to fetch sessions"
//	@Router			/chat/sessions [get]
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

		httputil.LogError(r, "failed to fetch chat sessions", err)

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

// GetSessionByID returns one chat session owned by the authenticated user.
//
//	@Summary		Get chat session by ID
//	@Description	Returns a single chat session, if it exists and is owned by the authenticated user.
//	@Tags			chat-sessions
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Chat session ID"
//	@Success		200	{object}	models.ChatSession
//	@Failure		400	{string}	string	"invalid session id"
//	@Failure		401	{string}	string	"unauthorized"
//	@Failure		404	{string}	string	"session not found"
//	@Failure		500	{string}	string	"failed to fetch session"
//	@Router			/chat/sessions/{id} [get]
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

		httputil.LogError(r, "failed to fetch chat session", err)

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

// UpdateSession replaces the title/mode/note_id of a chat session owned by
// the authenticated user.
//
//	@Summary		Update chat session
//	@Description	Replaces title, mode and note_id of an existing chat session owned by the authenticated user.
//	@Tags			chat-sessions
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Chat session ID"
//	@Param			session	body		models.ChatSession	true	"Updated title and mode"
//	@Success		200		{object}	models.ChatSession
//	@Failure		400		{string}	string	"invalid session id, title/mode missing, title too long, or mode is invalid"
//	@Failure		401		{string}	string	"unauthorized"
//	@Failure		404		{string}	string	"session not found"
//	@Failure		413		{string}	string	"request body too large"
//	@Failure		500		{string}	string	"failed to update session"
//	@Router			/chat/sessions/{id} [put]
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

	if err := httputil.DecodeJSON(w, r, httputil.MaxJSONBodyBytes, &session); err != nil {
		return
	}

	if err := validateChatSession(session); err != "" {
		http.Error(w, err, http.StatusBadRequest)
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

		httputil.LogError(r, "failed to update chat session", err)

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

// DeleteSession deletes a chat session (and its messages) owned by the
// authenticated user.
//
//	@Summary		Delete chat session
//	@Description	Deletes a chat session owned by the authenticated user.
//	@Tags			chat-sessions
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Chat session ID"
//	@Success		204	"No Content"
//	@Failure		400	{string}	string	"invalid session id"
//	@Failure		401	{string}	string	"unauthorized"
//	@Failure		404	{string}	string	"session not found"
//	@Failure		500	{string}	string	"failed to delete session"
//	@Router			/chat/sessions/{id} [delete]
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

		httputil.LogError(r, "failed to delete chat session", err)

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
