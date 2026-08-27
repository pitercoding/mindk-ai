package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/auth"
	"github.com/pitercoding/mindk-ai/backend/internal/httputil"
	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

const (
	maxNoteTitleLength   = 200
	maxNoteContentLength = 50_000
)

type NoteService interface {
	Create(note *models.Note) error
	GetAll(userID string) ([]models.Note, error)
	GetByID(id int, userID string) (*models.Note, error)
	Update(note *models.Note) error
	Delete(id int, userID string) error
}

type NoteHandler struct {
	Service NoteService
}

func NewNoteHandler(service NoteService) *NoteHandler {
	return &NoteHandler{
		Service: service,
	}
}

// CreateNote creates a note for the authenticated user.
//
//	@Summary		Create note
//	@Description	Creates a note owned by the authenticated user. user_id, id, created_at and updated_at are set by the server; any values sent by the client for them are ignored.
//	@Tags			notes
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			note	body		models.Note	true	"Note to create (title and content are read; other fields are ignored)"
//	@Success		201		{object}	models.Note
//	@Failure		400		{string}	string	"title and content are required, or exceed the maximum length"
//	@Failure		401		{string}	string	"unauthorized"
//	@Failure		413		{string}	string	"request body too large"
//	@Failure		500		{string}	string	"failed to create note"
//	@Router			/notes [post]
func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var note models.Note

	// 1. Read JSON from body
	if err := httputil.DecodeJSON(w, r, httputil.MaxJSONBodyBytes, &note); err != nil {
		return
	}

	// 2. Simple Validation
	if note.Title == "" || note.Content == "" {
		http.Error(w, "title and content are required", http.StatusBadRequest)
		return
	}

	if len(note.Title) > maxNoteTitleLength || len(note.Content) > maxNoteContentLength {
		http.Error(w, "title or content exceeds maximum length", http.StatusBadRequest)
		return
	}

	// 3. The owner is always the authenticated user, never the request body
	note.UserID = userID

	// 4. DB Saving
	err := h.Service.Create(&note)

	if err != nil {

		httputil.LogError(r, "failed to create note", err)

		http.Error(
			w,
			"failed to create note",
			http.StatusInternalServerError,
		)

		return
	}

	// 5. Response to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(note)
}

func (h *NoteHandler) HandleNotes(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		h.GetNotes(w, r)

	case http.MethodPost:
		h.CreateNote(w, r)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetNotes lists every note owned by the authenticated user.
//
//	@Summary		List notes
//	@Description	Returns every note owned by the authenticated user.
//	@Tags			notes
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.Note
//	@Failure		401	{string}	string	"unauthorized"
//	@Failure		500	{string}	string	"failed to fetch notes"
//	@Router			/notes [get]
func (h *NoteHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	notes, err := h.Service.GetAll(userID)
	if err != nil {
		httputil.LogError(r, "failed to fetch notes", err)
		http.Error(w, "failed to fetch notes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(notes)
}

// GetNoteByID returns one note owned by the authenticated user.
//
//	@Summary		Get note by ID
//	@Description	Returns a single note, if it exists and is owned by the authenticated user.
//	@Tags			notes
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Note ID"
//	@Success		200	{object}	models.Note
//	@Failure		400	{string}	string	"invalid note id"
//	@Failure		401	{string}	string	"unauthorized"
//	@Failure		404	{string}	string	"note not found"
//	@Failure		500	{string}	string	"failed to fetch note"
//	@Router			/notes/{id} [get]
func (h *NoteHandler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := httputil.GetIDFromPath(r)
	if err != nil {
		http.Error(w, "invalid note id", http.StatusBadRequest)
		return
	}

	note, err := h.Service.GetByID(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "note not found", http.StatusNotFound)
			return
		}

		httputil.LogError(r, "failed to fetch note", err)
		http.Error(w, "failed to fetch note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

func (h *NoteHandler) HandleNote(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		h.GetNoteByID(w, r)

	case http.MethodPut:
		h.UpdateNote(w, r)

	case http.MethodDelete:
		h.DeleteNote(w, r)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

// UpdateNote replaces the title and content of a note owned by the
// authenticated user.
//
//	@Summary		Update note
//	@Description	Replaces title and content of an existing note owned by the authenticated user. user_id is taken from the token, not the body.
//	@Tags			notes
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int			true	"Note ID"
//	@Param			note	body		models.Note	true	"Updated title and content"
//	@Success		200		{object}	models.Note
//	@Failure		400		{string}	string	"invalid note id, or title/content missing/too long"
//	@Failure		401		{string}	string	"unauthorized"
//	@Failure		404		{string}	string	"note not found"
//	@Failure		413		{string}	string	"request body too large"
//	@Failure		500		{string}	string	"failed to update note"
//	@Router			/notes/{id} [put]
func (h *NoteHandler) UpdateNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := httputil.GetIDFromPath(r)
	if err != nil {
		http.Error(w, "invalid note id", http.StatusBadRequest)
		return
	}

	var note models.Note

	if err := httputil.DecodeJSON(w, r, httputil.MaxJSONBodyBytes, &note); err != nil {
		return
	}

	if note.Title == "" || note.Content == "" {
		http.Error(w, "title and content are required", http.StatusBadRequest)
		return
	}

	if len(note.Title) > maxNoteTitleLength || len(note.Content) > maxNoteContentLength {
		http.Error(w, "title or content exceeds maximum length", http.StatusBadRequest)
		return
	}

	note.ID = id
	note.UserID = userID

	err = h.Service.Update(&note)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "note not found", http.StatusNotFound)
			return
		}

		httputil.LogError(r, "failed to update note", err)
		http.Error(w, "failed to update note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(note)
}

// DeleteNote deletes a note owned by the authenticated user.
//
//	@Summary		Delete note
//	@Description	Deletes a note owned by the authenticated user.
//	@Tags			notes
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Note ID"
//	@Success		204	"No Content"
//	@Failure		400	{string}	string	"invalid note id"
//	@Failure		401	{string}	string	"unauthorized"
//	@Failure		404	{string}	string	"note not found"
//	@Failure		500	{string}	string	"failed to delete note"
//	@Router			/notes/{id} [delete]
func (h *NoteHandler) DeleteNote(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := httputil.GetIDFromPath(r)
	if err != nil {
		http.Error(w, "invalid note id", http.StatusBadRequest)
		return
	}

	err = h.Service.Delete(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "note not found", http.StatusNotFound)
			return
		}

		httputil.LogError(r, "failed to delete note", err)
		http.Error(w, "failed to delete note", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
