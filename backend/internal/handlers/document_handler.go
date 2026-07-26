package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/httputil"
	"github.com/pitercoding/mindk-ai/backend/internal/models"
)

type DocumentService interface {
	Create(document *models.Document) error
	GetAll() ([]models.Document, error)
	GetByID(id int) (*models.Document, error)
	Delete(id int) error
	Search(query string) ([]models.Document, error)
}

type DocumentHandler struct {
	Service DocumentService
}

func NewDocumentHandler(service DocumentService) *DocumentHandler {
	return &DocumentHandler{Service: service}
}

func (h *DocumentHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	var document models.Document

	// 1. Read JSON from body
	err := json.NewDecoder(r.Body).Decode(&document)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	// 2. Simple Validation
	if document.Name == "" || document.Content == "" {
		http.Error(w, "Name and content are required", http.StatusBadRequest)
		return
	}

	// 3. DB Saving
	err = h.Service.Create(&document)

	if err != nil {

		http.Error(
			w,
			"failed to create document",
			http.StatusInternalServerError,
		)

		return
	}

	// 4. Response to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(document)
}

func (h *DocumentHandler) HandleDocuments(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		h.GetDocuments(w, r)

	case http.MethodPost:
		h.CreateDocument(w, r)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *DocumentHandler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	documents, err := h.Service.GetAll()
	if err != nil {
		http.Error(w, "failed to fetch documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(documents)
}

func (h *DocumentHandler) GetDocumentByID(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.GetIDFromPath(r)
	if err != nil {
		http.Error(w, "invalid document id", http.StatusBadRequest)
		return
	}

	document, err := h.Service.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "document not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to fetch document", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(document)
}

func (h *DocumentHandler) HandleDocument(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		h.GetDocumentByID(w, r)

	case http.MethodDelete:
		h.DeleteDocument(w, r)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *DocumentHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	id, err := httputil.GetIDFromPath(r)
	if err != nil {
		http.Error(w, "invalid document id", http.StatusBadRequest)
		return
	}

	err = h.Service.Delete(id)
	if err != nil {
		http.Error(w, "failed to delete document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
