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

const (
	// maxDocumentBytes bounds both the JSON create path (pasted content) and
	// the multipart upload path, since both end up as the same document
	// content and deserve the same generous ceiling.
	maxDocumentBytes      = 10 << 20 // 10 MiB
	maxDocumentNameLength = 255
	maxSearchQueryLength  = 500
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

// CreateDocument creates a document from pasted/JSON text content, owned by
// the authenticated user.
//
//	@Summary		Create document from text
//	@Description	Creates a document from JSON content (as opposed to a file upload). user_id, id and created_at are set by the server.
//	@Tags			documents
//	@Security		BearerAuth
//	@Accept			json
//	@Produce		json
//	@Param			document	body		models.Document	true	"Document to create (name and content are read; other fields are ignored)"
//	@Success		201			{object}	models.Document
//	@Failure		400			{string}	string	"name and content are required, or name exceeds maximum length"
//	@Failure		401			{string}	string	"unauthorized"
//	@Failure		413			{string}	string	"request body too large"
//	@Failure		500			{string}	string	"failed to create document"
//	@Router			/documents [post]
func (h *DocumentHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var document models.Document

	// 1. Read JSON from body
	if err := httputil.DecodeJSON(w, r, maxDocumentBytes, &document); err != nil {
		return
	}

	// 2. Simple Validation
	if document.Name == "" || document.Content == "" {
		http.Error(w, "Name and content are required", http.StatusBadRequest)
		return
	}

	if len(document.Name) > maxDocumentNameLength {
		http.Error(w, "name exceeds maximum length", http.StatusBadRequest)
		return
	}

	// 3. The owner is always the authenticated user, never the request body
	document.UserID = userID

	// 4. DB Saving
	err := h.Service.Create(&document)

	if err != nil {

		httputil.LogError(r, "failed to create document", err)

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

// UploadDocument creates a document from an uploaded file, owned by the
// authenticated user.
//
//	@Summary		Upload document file
//	@Description	Creates a document from an uploaded file. The file's text is extracted server-side (by extension) and stored as the document content. Maximum upload size is 10 MiB.
//	@Tags			documents
//	@Security		BearerAuth
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file	true	"File to upload"
//	@Success		201		{object}	models.Document
//	@Failure		400		{string}	string	"file is required, or invalid multipart form"
//	@Failure		401		{string}	string	"unauthorized"
//	@Failure		413		{string}	string	"request body too large (over 10 MiB)"
//	@Failure		500		{string}	string	"failed to read/save file"
//	@Router			/documents/upload [post]
func (h *DocumentHandler) UploadDocument(
	w http.ResponseWriter,
	r *http.Request,
) {

	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxDocumentBytes)

	err := r.ParseMultipartForm(maxDocumentBytes)
	if err != nil {

		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			return
		}

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

		httputil.LogError(r, "failed to read uploaded file", err)

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

		httputil.LogError(r, "failed to save document", err)

		http.Error(
			w,
			"failed to save document",
			http.StatusInternalServerError,
		)

		return
	}

	// filename is deliberately omitted: it is user-supplied and may carry
	// personal information, unlike the size/type of the upload.
	httputil.LogInfo(r, "document uploaded", "size_bytes", len(content), "extension", extension)

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

// GetDocuments lists every document owned by the authenticated user.
//
//	@Summary		List documents
//	@Description	Returns every document owned by the authenticated user.
//	@Tags			documents
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.Document
//	@Failure		401	{string}	string	"unauthorized"
//	@Failure		500	{string}	string	"failed to fetch documents"
//	@Router			/documents [get]
func (h *DocumentHandler) GetDocuments(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	documents, err := h.Service.GetAll(userID)
	if err != nil {
		httputil.LogError(r, "failed to fetch documents", err)
		http.Error(w, "failed to fetch documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(documents)
}

// GetDocumentByID returns one document owned by the authenticated user.
//
//	@Summary		Get document by ID
//	@Description	Returns a single document, if it exists and is owned by the authenticated user.
//	@Tags			documents
//	@Security		BearerAuth
//	@Produce		json
//	@Param			id	path		int	true	"Document ID"
//	@Success		200	{object}	models.Document
//	@Failure		400	{string}	string	"invalid document id"
//	@Failure		401	{string}	string	"unauthorized"
//	@Failure		404	{string}	string	"document not found"
//	@Failure		500	{string}	string	"failed to fetch document"
//	@Router			/documents/{id} [get]
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

		httputil.LogError(r, "failed to fetch document", err)
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

// DeleteDocument deletes a document owned by the authenticated user.
//
//	@Summary		Delete document
//	@Description	Deletes a document (and its derived chunks/embeddings) owned by the authenticated user.
//	@Tags			documents
//	@Security		BearerAuth
//	@Param			id	path	int	true	"Document ID"
//	@Success		204	"No Content"
//	@Failure		400	{string}	string	"invalid document id"
//	@Failure		401	{string}	string	"unauthorized"
//	@Failure		404	{string}	string	"document not found"
//	@Failure		500	{string}	string	"failed to delete document"
//	@Router			/documents/{id} [delete]
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

		httputil.LogError(r, "failed to delete document", err)
		http.Error(w, "failed to delete document", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// SearchDocuments performs a semantic search over the authenticated user's
// documents and returns the matching documents.
//
//	@Summary		Search documents
//	@Description	Searches the authenticated user's documents by semantic relevance to the query.
//	@Tags			documents
//	@Security		BearerAuth
//	@Produce		json
//	@Param			q	query		string	true	"Search query (max 500 characters)"
//	@Success		200	{array}		models.Document
//	@Failure		400	{string}	string	"query parameter q is required, or exceeds maximum length"
//	@Failure		401	{string}	string	"unauthorized"
//	@Failure		500	{string}	string	"failed to search documents"
//	@Router			/documents/search [get]
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

	if query == "" {
		http.Error(w, "query parameter q is required", http.StatusBadRequest)
		return
	}

	if len(query) > maxSearchQueryLength {
		http.Error(w, "query exceeds maximum length", http.StatusBadRequest)
		return
	}

	documents, err := h.Service.Search(query, userID)
	if err != nil {

		// The query itself is not logged: it is free-text user input and may
		// contain the same kind of personal content documents do.
		httputil.LogError(r, "failed to search documents", err)

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
