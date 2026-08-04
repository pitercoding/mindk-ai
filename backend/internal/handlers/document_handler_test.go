package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/handlers/mocks"
	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDocumentHandlerCreateDocument(t *testing.T) {

	service := &mocks.FakeDocumentService{}

	handler := NewDocumentHandler(service)

	body := `
	{
		"name": "test.md",
		"type": ".md",
		"content": "document content"
	}
	`

	request := httptest.NewRequest(
		http.MethodPost,
		"/documents",
		bytes.NewBufferString(body),
	)

	request.Header.Set(
		"Content-Type",
		"application/json",
	)

	response := httptest.NewRecorder()

	handler.CreateDocument(
		response,
		request,
	)

	require.Equal(
		t,
		http.StatusCreated,
		response.Code,
	)

	assert.True(
		t,
		service.Created,
	)

	assert.Equal(
		t,
		"test.md",
		service.Document.Name,
	)
}

func TestDocumentHandlerCreateDocument_InvalidJSON(t *testing.T) {

	service := &mocks.FakeDocumentService{}

	handler := NewDocumentHandler(service)

	request := httptest.NewRequest(
		http.MethodPost,
		"/documents",
		bytes.NewBufferString("{invalid"),
	)

	response := httptest.NewRecorder()

	handler.CreateDocument(
		response,
		request,
	)

	assert.Equal(
		t,
		http.StatusBadRequest,
		response.Code,
	)
}

func TestDocumentHandlerCreateDocument_ServiceError(t *testing.T) {

	service := &mocks.FakeDocumentService{
		Err: errors.New("database error"),
	}

	handler := NewDocumentHandler(service)

	body := `
	{
		"name": "test.md",
		"content": "content"
	}
	`

	request := httptest.NewRequest(
		http.MethodPost,
		"/documents",
		bytes.NewBufferString(body),
	)

	response := httptest.NewRecorder()

	handler.CreateDocument(
		response,
		request,
	)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		response.Code,
	)
}

func TestDocumentHandlerGetDocuments(t *testing.T) {

	service := &mocks.FakeDocumentService{
		Documents: []models.Document{
			{
				ID:      1,
				Name:    "test.md",
				Content: "content",
			},
		},
	}

	handler := NewDocumentHandler(service)

	request := httptest.NewRequest(
		http.MethodGet,
		"/documents",
		nil,
	)

	response := httptest.NewRecorder()

	handler.GetDocuments(
		response,
		request,
	)

	require.Equal(
		t,
		http.StatusOK,
		response.Code,
	)

	var documents []models.Document

	err := json.NewDecoder(
		response.Body,
	).Decode(&documents)

	require.NoError(
		t,
		err,
	)

	assert.Len(
		t,
		documents,
		1,
	)
}

func TestDocumentHandlerDeleteDocument(t *testing.T) {

	service := &mocks.FakeDocumentService{}

	handler := NewDocumentHandler(service)

	request := httptest.NewRequest(
		http.MethodDelete,
		"/documents/1",
		nil,
	)

	response := httptest.NewRecorder()

	handler.DeleteDocument(
		response,
		request,
	)

	assert.Equal(
		t,
		http.StatusNoContent,
		response.Code,
	)

	assert.True(
		t,
		service.Deleted,
	)
}

func TestDocumentHandlerUploadDocument(t *testing.T) {

	service := &mocks.FakeDocumentService{}

	handler := NewDocumentHandler(service)

	var body bytes.Buffer

	writer := multipart.NewWriter(&body)

	fileWriter, err := writer.CreateFormFile(
		"file",
		"test.txt",
	)

	require.NoError(
		t,
		err,
	)

	_, err = fileWriter.Write(
		[]byte("This is a test document content."),
	)

	require.NoError(
		t,
		err,
	)

	err = writer.Close()

	require.NoError(
		t,
		err,
	)

	request := httptest.NewRequest(
		http.MethodPost,
		"/documents/upload",
		&body,
	)

	request.Header.Set(
		"Content-Type",
		writer.FormDataContentType(),
	)

	response := httptest.NewRecorder()

	handler.UploadDocument(
		response,
		request,
	)

	require.Equal(
		t,
		http.StatusCreated,
		response.Code,
	)

	assert.True(
		t,
		service.Created,
	)

	assert.Equal(
		t,
		"test.txt",
		service.Document.Name,
	)

	assert.Equal(
		t,
		".txt",
		service.Document.Type,
	)

	assert.Equal(
		t,
		"This is a test document content.",
		service.Document.Content,
	)
}
