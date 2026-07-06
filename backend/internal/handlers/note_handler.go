package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
)

type NoteHandler struct {
	Repo *repository.NoteRepository
}

func NewNoteHandler(repo *repository.NoteRepository) *NoteHandler {
	return &NoteHandler{Repo: repo}
}

func (h *NoteHandler) CreateNote(w http.ResponseWriter, r *http.Request) {
	var note models.Note

	// 1. Read JSON from body
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Simple Validation
	if note.Title == "" || note.Content == "" {
		http.Error(w, "title and content are required", http.StatusBadRequest)
		return
	}

	// 3. DB Saving
	err = h.Repo.Create(&note)
	if err != nil {
		http.Error(w, "failed to create note", http.StatusInternalServerError)
		return
	}

	// 4. Response to JSON
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
	notes, err := h.Repo.GetAll()
	if err != nil {
		http.Error(w, "failed to fetch notes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(notes)
}

func (h *NoteHandler) GetNoteByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	idStr := parts[len(parts)-1]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid note id", http.StatusBadRequest)
		return
	}

	note, err := h.Repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "note not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to fetch note", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(note)
}
