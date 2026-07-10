package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/handlers/mocks"
	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// GET //
func TestNoteHandlerGetNotes(t *testing.T) {

	service := &mocks.FakeNoteService{
		Notes: []models.Note{
			{
				ID:      1,
				Title:   "Go",
				Content: "Go is awesome",
			},
			{
				ID:      2,
				Title:   "Docker",
				Content: "Docker containers",
			},
		},
	}

	handler := NewNoteHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/notes",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetNotes(
		recorder,
		req,
	)

	require.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	var response []models.Note

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&response)

	require.NoError(t, err)

	assert.Len(t, response, 2)
	assert.Equal(t, "Go", response[0].Title)
	assert.Equal(t, "Docker", response[1].Title)
}

func TestNoteHandlerGetNotes_ServiceError(t *testing.T) {

	service := &mocks.FakeNoteService{
		Err: errors.New("database error"),
	}

	handler := NewNoteHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/notes",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetNotes(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		recorder.Code,
	)

	assert.Contains(
		t,
		recorder.Body.String(),
		"failed to fetch notes",
	)
}

// CREATE //
func TestNoteHandlerCreateNote(t *testing.T) {

	service := &mocks.FakeNoteService{}

	handler := NewNoteHandler(service)

	body := `{
		"title":"Go",
		"content":"Go is awesome"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/notes",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.CreateNote(
		recorder,
		req,
	)

	require.Equal(
		t,
		http.StatusCreated,
		recorder.Code,
	)

	var response models.Note

	err := json.NewDecoder(recorder.Body).Decode(&response)

	require.NoError(t, err)

	assert.Equal(t, "Go", response.Title)
	assert.Equal(t, "Go is awesome", response.Content)

	require.NotNil(t, service.CreatedNote)

	assert.Equal(t, "Go", service.CreatedNote.Title)
	assert.Equal(t, "Go is awesome", service.CreatedNote.Content)
}

func TestNoteHandlerCreateNote_InvalidJSON(t *testing.T) {

	service := &mocks.FakeNoteService{}

	handler := NewNoteHandler(service)

	body := `{"title":"Go"`

	req := httptest.NewRequest(
		http.MethodPost,
		"/notes",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.CreateNote(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusBadRequest,
		recorder.Code,
	)

	assert.Contains(
		t,
		recorder.Body.String(),
		"invalid request body",
	)

	assert.Nil(t, service.CreatedNote)
}

func TestNoteHandlerCreateNote_Validation(t *testing.T) {

	tests := []struct {
		name string
		body string
	}{
		{
			name: "missing title",
			body: `{
				"content":"Go is awesome"
			}`,
		},
		{
			name: "missing content",
			body: `{
				"title":"Go"
			}`,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			service := &mocks.FakeNoteService{}

			handler := NewNoteHandler(service)

			req := httptest.NewRequest(
				http.MethodPost,
				"/notes",
				strings.NewReader(tt.body),
			)

			req.Header.Set(
				"Content-Type",
				"application/json",
			)

			recorder := httptest.NewRecorder()

			handler.CreateNote(
				recorder,
				req,
			)

			assert.Equal(
				t,
				http.StatusBadRequest,
				recorder.Code,
			)

			assert.Contains(
				t,
				recorder.Body.String(),
				"title and content are required",
			)

			assert.Nil(t, service.CreatedNote)
		})
	}
}

func TestNoteHandlerCreateNote_ServiceError(t *testing.T) {

	service := &mocks.FakeNoteService{
		Err: errors.New("database error"),
	}

	handler := NewNoteHandler(service)

	body := `{
		"title":"Go",
		"content":"Go is awesome"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/notes",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.CreateNote(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		recorder.Code,
	)

	assert.Contains(
		t,
		recorder.Body.String(),
		"failed to create note",
	)

	require.NotNil(t, service.CreatedNote)

	assert.Equal(
		t,
		"Go",
		service.CreatedNote.Title,
	)
}

// GET BY ID //
func TestNoteHandlerGetNoteByID(t *testing.T) {

	service := &mocks.FakeNoteService{
		Note: &models.Note{
			ID:      1,
			Title:   "Go",
			Content: "Go is awesome",
		},
	}

	handler := NewNoteHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/notes/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetNoteByID(
		recorder,
		req,
	)

	require.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	var response models.Note

	err := json.NewDecoder(recorder.Body).Decode(&response)

	require.NoError(t, err)

	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "Go", response.Title)
	assert.Equal(t, "Go is awesome", response.Content)
}

func TestNoteHandlerGetNoteByID_InvalidID(t *testing.T) {

	service := &mocks.FakeNoteService{}

	handler := NewNoteHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/notes/abc",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetNoteByID(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusBadRequest,
		recorder.Code,
	)

	assert.Contains(
		t,
		recorder.Body.String(),
		"invalid note id",
	)
}

func TestNoteHandlerGetNoteByID_NotFound(t *testing.T) {

	service := &mocks.FakeNoteService{
		Err: sql.ErrNoRows,
	}

	handler := NewNoteHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/notes/99",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetNoteByID(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusNotFound,
		recorder.Code,
	)

	assert.Contains(
		t,
		recorder.Body.String(),
		"note not found",
	)
}

func TestNoteHandlerGetNoteByID_ServiceError(t *testing.T) {

	service := &mocks.FakeNoteService{
		Err: errors.New("database error"),
	}

	handler := NewNoteHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/notes/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetNoteByID(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		recorder.Code,
	)

	assert.Contains(
		t,
		recorder.Body.String(),
		"failed to fetch note",
	)
}

// UPDATE //
func TestNoteHandlerUpdateNote(t *testing.T) {

	service := &mocks.FakeNoteService{}

	handler := NewNoteHandler(service)

	body := `{
		"title":"Go Updated",
		"content":"Go is awesome and fast"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/notes/1",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.UpdateNote(
		recorder,
		req,
	)

	require.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	var response models.Note

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&response)

	require.NoError(t, err)

	assert.Equal(t, 1, response.ID)
	assert.Equal(t, "Go Updated", response.Title)
	assert.Equal(t, "Go is awesome and fast", response.Content)

	require.NotNil(t, service.UpdatedNote)

	assert.Equal(t, 1, service.UpdatedNote.ID)
	assert.Equal(t, "Go Updated", service.UpdatedNote.Title)
}

func TestNoteHandlerUpdateNote_InvalidID(t *testing.T) {

	service := &mocks.FakeNoteService{}

	handler := NewNoteHandler(service)

	body := `{
		"title":"Go Updated",
		"content":"Go is awesome and fast"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/notes/abc",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.UpdateNote(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusBadRequest,
		recorder.Code,
	)

	assert.Contains(
		t,
		recorder.Body.String(),
		"invalid note id",
	)

	assert.Nil(
		t,
		service.UpdatedNote,
	)
}

func TestNoteHandlerUpdateNote_InvalidJSON(t *testing.T) {

	service := &mocks.FakeNoteService{}

	handler := NewNoteHandler(service)

	body := `{
		"title":"Go Updated"
	`

	req := httptest.NewRequest(
		http.MethodPut,
		"/notes/1",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.UpdateNote(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusBadRequest,
		recorder.Code,
	)

	assert.Contains(
		t,
		recorder.Body.String(),
		"invalid request body",
	)

	assert.Nil(
		t,
		service.UpdatedNote,
	)
}

func TestNoteHandlerUpdateNote_ServiceError(t *testing.T) {

	service := &mocks.FakeNoteService{
		Err: errors.New("database error"),
	}

	handler := NewNoteHandler(service)

	body := `{
		"title":"Go Updated",
		"content":"Go is awesome and fast"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/notes/1",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.UpdateNote(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		recorder.Code,
	)

	assert.Contains(
		t,
		recorder.Body.String(),
		"failed to update note",
	)

	require.NotNil(
		t,
		service.UpdatedNote,
	)

	assert.Equal(
		t,
		1,
		service.UpdatedNote.ID,
	)
}

// DELETE //
func TestNoteHandlerDeleteNote(t *testing.T) {

	service := &mocks.FakeNoteService{}

	handler := NewNoteHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/notes/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.DeleteNote(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusNoContent,
		recorder.Code,
	)

	assert.Equal(
		t,
		1,
		service.DeletedID,
	)
}

func TestNoteHandlerDeleteNote_InvalidID(t *testing.T) {

	service := &mocks.FakeNoteService{}

	handler := NewNoteHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/notes/abc",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.DeleteNote(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusBadRequest,
		recorder.Code,
	)

	assert.Contains(
		t,
		recorder.Body.String(),
		"invalid note id",
	)

	assert.Equal(
		t,
		0,
		service.DeletedID,
	)
}

func TestNoteHandlerDeleteNote_ServiceError(t *testing.T) {

	service := &mocks.FakeNoteService{
		Err: errors.New("database error"),
	}

	handler := NewNoteHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/notes/1",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.DeleteNote(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		recorder.Code,
	)

	assert.Contains(
		t,
		recorder.Body.String(),
		"failed to delete note",
	)

	assert.Equal(
		t,
		1,
		service.DeletedID,
	)
}
