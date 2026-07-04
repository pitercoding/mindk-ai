package handlers

import (
	"encoding/json"
	"net/http"

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
