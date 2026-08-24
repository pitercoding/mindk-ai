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
