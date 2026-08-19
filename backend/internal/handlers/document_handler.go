package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"path/filepath"

	"github.com/pitercoding/mindk-ai/backend/internal/auth"
	"github.com/pitercoding/mindk-ai/backend/internal/httputil"
	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/utils"
)

type DocumentService interface {
	Create(document *models.Document) error
	GetAll(userID string) ([]models.Document, error)
	GetByID(id int, userID string) (*models.Document, error)
	Delete(id int, userID string) error
	Search(query string, userID string) ([]models.Document, error)
}

type DocumentHandler struct {
	Service DocumentService
}

func NewDocumentHandler(service DocumentService) *DocumentHandler {
	return &DocumentHandler{Service: service}
}

func (h *DocumentHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

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

	// 3. The owner is always the authenticated user, never the request body
	document.UserID = userID

	// 4. DB Saving
	err = h.Service.Create(&document)

	if err != nil {

		http.Error(
			w,
			"failed to create document",
			http.StatusInternalServerError,
		)

		return
	}

	// 5. Response to JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(document)
}

func (h *DocumentHandler) UploadDocument(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {

		http.Error(
			w,
			"invalid multipart form",
			http.StatusBadRequest,
		)

		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {

		http.Error(
			w,
			"file is required",
			http.StatusBadRequest,
		)

		return
	}

	defer file.Close()

	extension := filepath.Ext(header.Filename)

	content, err := utils.ReadFile(
		file,
		extension,
	)
	if err != nil {

		http.Error(
			w,
			"failed to read file",
			http.StatusInternalServerError,
		)

		return
	}

	document := models.Document{
		UserID:  userID,
		Name:    header.Filename,
		Type:    extension,
		Content: content,
	}

	err = h.Service.Create(&document)
	if err != nil {

		http.Error(
			w,
			"failed to save document",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

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
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	documents, err := h.Service.GetAll(userID)
	if err != nil {
		http.Error(w, "failed to fetch documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(documents)
}

func (h *DocumentHandler) GetDocumentByID(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := httputil.GetIDFromPath(r)
	if err != nil {
		http.Error(w, "invalid document id", http.StatusBadRequest)
		return
	}

	document, err := h.Service.GetByID(id, userID)
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
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := httputil.GetIDFromPath(r)
	if err != nil {
		http.Error(w, "invalid document id", http.StatusBadRequest)
		return
	}

	err = h.Service.Delete(id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "document not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to delete document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *DocumentHandler) SearchDocuments(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	query := r.URL.Query().Get("q")

	documents, err := h.Service.Search(query, userID)
	if err != nil {

		http.Error(
			w,
			"failed to search documents",
			http.StatusInternalServerError,
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(documents)
}
