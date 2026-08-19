package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
	"github.com/pitercoding/mindk-ai/backend/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatSessionHandlerCreateSession(t *testing.T) {
	service := &mocks.FakeChatSessionService{
		Session: &models.ChatSession{
			ID:    1,
			Title: "Docker",
			Mode:  "knowledge",
		},
	}

	handler := NewChatSessionHandler(service)

	body := `{
		"title": "Docker",
		"mode": "knowledge"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat/sessions",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSessions(
		recorder,
		req,
	)

	require.Equal(
		t,
		http.StatusCreated,
		recorder.Code,
	)

	var response models.ChatSession

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&response)

	require.NoError(t, err)

	assert.Equal(
		t,
		1,
		response.ID,
	)

	assert.Equal(
		t,
		"Docker",
		response.Title,
	)

	assert.Equal(
		t,
		"knowledge",
		response.Mode,
	)

	require.Len(
		t,
		service.Created,
		1,
	)

	assert.Equal(
		t,
		"Docker",
		service.Created[0].Title,
	)
}

func TestChatSessionHandlerCreateSession_IgnoresBodyUserID(t *testing.T) {
	service := &mocks.FakeChatSessionService{
		Session: &models.ChatSession{ID: 1},
	}

	handler := NewChatSessionHandler(service)

	body := `{
		"user_id": "someone_else",
		"title": "Docker",
		"mode": "knowledge"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat/sessions",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSessions(
		recorder,
		req,
	)

	require.Equal(
		t,
		http.StatusCreated,
		recorder.Code,
	)

	require.Len(t, service.Created, 1)
	assert.Equal(t, testUserID, service.Created[0].UserID)
}

func TestChatSessionHandlerCreateSession_Unauthorized(t *testing.T) {
	service := &mocks.FakeChatSessionService{}

	handler := NewChatSessionHandler(service)

	body := `{
		"title": "Docker",
		"mode": "knowledge"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat/sessions",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.HandleSessions(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusUnauthorized,
		recorder.Code,
	)

	assert.Empty(t, service.Created)
}

func TestChatSessionHandlerCreateSessionInvalidJSON(t *testing.T) {
	service := &mocks.FakeChatSessionService{}

	handler := NewChatSessionHandler(service)

	body := `{
		"title": "Docker",
		"mode": "knowledge"
	`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat/sessions",
		strings.NewReader(body),
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSessions(
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

	assert.Empty(
		t,
		service.Created,
	)
}

func TestChatSessionHandlerCreateSessionMissingFields(t *testing.T) {
	service := &mocks.FakeChatSessionService{}

	handler := NewChatSessionHandler(service)

	body := `{
		"title": "Docker"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat/sessions",
		strings.NewReader(body),
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSessions(
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
		"title and mode are required",
	)

	assert.Empty(
		t,
		service.Created,
	)
}

func TestChatSessionHandlerCreateSessionServiceError(t *testing.T) {
	service := &mocks.FakeChatSessionService{
		Err: errors.New("database error"),
	}

	handler := NewChatSessionHandler(service)

	body := `{
		"title": "Docker",
		"mode": "knowledge"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat/sessions",
		strings.NewReader(body),
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSessions(
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
		"failed to create session",
	)

	assert.Len(
		t,
		service.Created,
		1,
	)
}

func TestChatSessionHandlerGetAllSessions(t *testing.T) {
	service := &mocks.FakeChatSessionService{
		Sessions: []models.ChatSession{
			{
				ID:    1,
				Title: "Docker",
				Mode:  "knowledge",
			},
			{
				ID:    2,
				Title: "Go",
				Mode:  "note",
			},
		},
	}

	handler := NewChatSessionHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/chat/sessions",
		nil,
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSessions(
		recorder,
		req,
	)

	require.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	var response []models.ChatSession

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&response)

	require.NoError(t, err)

	require.Len(
		t,
		response,
		2,
	)

	assert.Equal(
		t,
		"Docker",
		response[0].Title,
	)

	assert.Equal(
		t,
		"Go",
		response[1].Title,
	)
}

func TestChatSessionHandlerGetAllSessionsServiceError(t *testing.T) {
	service := &mocks.FakeChatSessionService{
		Err: errors.New("database error"),
	}

	handler := NewChatSessionHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/chat/sessions",
		nil,
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSessions(
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
		"failed to fetch sessions",
	)
}

func TestChatSessionHandlerGetAllSessions_Unauthorized(t *testing.T) {
	service := &mocks.FakeChatSessionService{}

	handler := NewChatSessionHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/chat/sessions",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.HandleSessions(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusUnauthorized,
		recorder.Code,
	)
}

func TestChatSessionHandlerGetSessionByID(t *testing.T) {
	service := &mocks.FakeChatSessionService{
		Session: &models.ChatSession{
			ID:    5,
			Title: "Go",
			Mode:  "note",
		},
	}

	handler := NewChatSessionHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/chat/sessions/5",
		nil,
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSession(
		recorder,
		req,
	)

	require.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	var response models.ChatSession

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&response)

	require.NoError(t, err)

	assert.Equal(
		t,
		5,
		response.ID,
	)

	assert.Equal(
		t,
		"Go",
		response.Title,
	)

	assert.Equal(
		t,
		"note",
		response.Mode,
	)
}

func TestChatSessionHandlerGetSessionByIDNotFound(t *testing.T) {
	service := &mocks.FakeChatSessionService{
		Err: repository.ErrChatSessionNotFound,
	}

	handler := NewChatSessionHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/chat/sessions/999",
		nil,
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSession(
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
		"session not found",
	)
}

func TestChatSessionHandlerGetSessionByIDInvalidID(t *testing.T) {
	service := &mocks.FakeChatSessionService{}

	handler := NewChatSessionHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/chat/sessions/abc",
		nil,
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSession(
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
		"invalid session id",
	)
}

func TestChatSessionHandlerGetSessionByID_Unauthorized(t *testing.T) {
	service := &mocks.FakeChatSessionService{}

	handler := NewChatSessionHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/chat/sessions/5",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.HandleSession(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusUnauthorized,
		recorder.Code,
	)
}

func TestChatSessionHandlerUpdateSession(t *testing.T) {
	service := &mocks.FakeChatSessionService{}

	handler := NewChatSessionHandler(service)

	body := `{
		"title": "Updated title",
		"mode": "knowledge"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/chat/sessions/10",
		strings.NewReader(body),
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSession(
		recorder,
		req,
	)

	require.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	require.Len(
		t,
		service.Updated,
		1,
	)

	assert.Equal(
		t,
		10,
		service.Updated[0].ID,
	)

	assert.Equal(
		t,
		"Updated title",
		service.Updated[0].Title,
	)

	assert.Equal(
		t,
		"knowledge",
		service.Updated[0].Mode,
	)

	var response models.ChatSession

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&response)

	require.NoError(t, err)

	assert.Equal(
		t,
		10,
		response.ID,
	)
}

func TestChatSessionHandlerUpdateSession_IgnoresBodyUserID(t *testing.T) {
	service := &mocks.FakeChatSessionService{}

	handler := NewChatSessionHandler(service)

	body := `{
		"user_id": "someone_else",
		"title": "Updated title",
		"mode": "knowledge"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/chat/sessions/10",
		strings.NewReader(body),
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSession(
		recorder,
		req,
	)

	require.Equal(t, http.StatusOK, recorder.Code)
	require.Len(t, service.Updated, 1)
	assert.Equal(t, testUserID, service.Updated[0].UserID)
}

func TestChatSessionHandlerUpdateSession_Unauthorized(t *testing.T) {
	service := &mocks.FakeChatSessionService{}

	handler := NewChatSessionHandler(service)

	body := `{
		"title": "Updated title",
		"mode": "knowledge"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/chat/sessions/10",
		strings.NewReader(body),
	)

	recorder := httptest.NewRecorder()

	handler.HandleSession(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusUnauthorized,
		recorder.Code,
	)

	assert.Empty(t, service.Updated)
}

func TestChatSessionHandlerUpdateSessionInvalidID(t *testing.T) {
	service := &mocks.FakeChatSessionService{}

	handler := NewChatSessionHandler(service)

	body := `{
		"title": "Updated title",
		"mode": "knowledge"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/chat/sessions/abc",
		strings.NewReader(body),
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSession(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusBadRequest,
		recorder.Code,
	)

	assert.Empty(
		t,
		service.Updated,
	)
}

func TestChatSessionHandlerUpdateSessionInvalidJSON(t *testing.T) {
	service := &mocks.FakeChatSessionService{}

	handler := NewChatSessionHandler(service)

	body := `{
		"title": "Updated title",
		"mode": "knowledge"
	`

	req := httptest.NewRequest(
		http.MethodPut,
		"/chat/sessions/10",
		strings.NewReader(body),
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSession(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusBadRequest,
		recorder.Code,
	)

	assert.Empty(
		t,
		service.Updated,
	)
}

func TestChatSessionHandlerUpdateSessionMissingFields(t *testing.T) {
	service := &mocks.FakeChatSessionService{}

	handler := NewChatSessionHandler(service)

	body := `{
		"title": "Updated title"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/chat/sessions/10",
		strings.NewReader(body),
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSession(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusBadRequest,
		recorder.Code,
	)

	assert.Empty(
		t,
		service.Updated,
	)
}

func TestChatSessionHandlerUpdateSessionNotFound(t *testing.T) {
	service := &mocks.FakeChatSessionService{
		Err: repository.ErrChatSessionNotFound,
	}

	handler := NewChatSessionHandler(service)

	body := `{
		"title": "Updated title",
		"mode": "knowledge"
	}`

	req := httptest.NewRequest(
		http.MethodPut,
		"/chat/sessions/10",
		strings.NewReader(body),
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSession(
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
		"session not found",
	)
}

func TestChatSessionHandlerDeleteSession(t *testing.T) {
	service := &mocks.FakeChatSessionService{}

	handler := NewChatSessionHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/chat/sessions/10",
		nil,
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSession(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusNoContent,
		recorder.Code,
	)

	require.Len(
		t,
		service.Deleted,
		1,
	)

	assert.Equal(
		t,
		10,
		service.Deleted[0],
	)
}

func TestChatSessionHandlerDeleteSession_Unauthorized(t *testing.T) {
	service := &mocks.FakeChatSessionService{}

	handler := NewChatSessionHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/chat/sessions/10",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.HandleSession(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusUnauthorized,
		recorder.Code,
	)

	assert.Empty(t, service.Deleted)
}

func TestChatSessionHandlerDeleteSessionInvalidID(t *testing.T) {
	service := &mocks.FakeChatSessionService{}

	handler := NewChatSessionHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/chat/sessions/abc",
		nil,
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSession(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusBadRequest,
		recorder.Code,
	)

	assert.Empty(
		t,
		service.Deleted,
	)
}

func TestChatSessionHandlerDeleteSessionNotFound(t *testing.T) {
	service := &mocks.FakeChatSessionService{
		Err: repository.ErrChatSessionNotFound,
	}

	handler := NewChatSessionHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/chat/sessions/10",
		nil,
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSession(
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
		"session not found",
	)
}

func TestChatSessionHandlerDeleteSessionServiceError(t *testing.T) {
	service := &mocks.FakeChatSessionService{
		Err: errors.New("database error"),
	}

	handler := NewChatSessionHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/chat/sessions/10",
		nil,
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.HandleSession(
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
		"failed to delete session",
	)
}

func TestChatSessionHandlerMethodsNotAllowed(t *testing.T) {
	service := &mocks.FakeChatSessionService{}

	handler := NewChatSessionHandler(service)

	t.Run("collection", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodDelete,
			"/chat/sessions",
			nil,
		)

		recorder := httptest.NewRecorder()

		handler.HandleSessions(
			recorder,
			req,
		)

		assert.Equal(
			t,
			http.StatusMethodNotAllowed,
			recorder.Code,
		)
	})

	t.Run("single session", func(t *testing.T) {
		req := httptest.NewRequest(
			http.MethodPost,
			"/chat/sessions/10",
			nil,
		)

		recorder := httptest.NewRecorder()

		handler.HandleSession(
			recorder,
			req,
		)

		assert.Equal(
			t,
			http.StatusMethodNotAllowed,
			recorder.Code,
		)
	})
}
