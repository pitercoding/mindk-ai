package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/handlers/mocks"
	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatHandlerAsk(t *testing.T) {
	service := &mocks.FakeChatService{
		Answer:          "Docker is a container platform.",
		ResponseSession: 1,
	}

	handler := NewChatHandler(service)

	body := `{
	"session_id": 1,
	"message": "What do my notes say about Docker?"
}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.Ask(
		recorder,
		req,
	)

	require.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	var response models.ChatResponse

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&response)

	require.NoError(t, err)

	assert.Equal(
		t,
		"Docker is a container platform.",
		response.Answer,
	)

	assert.Equal(
		t,
		1,
		response.SessionID,
	)

	assert.True(
		t,
		service.Called,
	)

	assert.Equal(
		t,
		1,
		service.LastSessionID,
	)

	assert.Equal(
		t,
		"What do my notes say about Docker?",
		service.LastMessage,
	)

}

func TestChatHandlerAsk_AutoCreatesSession(t *testing.T) {
	service := &mocks.FakeChatService{
		Answer:          "Go is awesome.",
		ResponseSession: 42,
	}

	handler := NewChatHandler(service)

	body := `{
	"message": "What do my notes say about Go?",
	"mode": "knowledge",
	"title": "My chat"
}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.Ask(
		recorder,
		req,
	)

	require.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	var response models.ChatResponse

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&response)

	require.NoError(t, err)

	assert.Equal(
		t,
		42,
		response.SessionID,
	)

	assert.Equal(
		t,
		0,
		service.LastSessionID,
	)

	assert.Equal(
		t,
		"knowledge",
		service.LastMode,
	)

	assert.Equal(
		t,
		"My chat",
		service.LastTitle,
	)
}

func TestChatHandlerAsk_InvalidJSON(t *testing.T) {
	service := &mocks.FakeChatService{}

	handler := NewChatHandler(service)

	body := `{
	"session_id": 1,
	"message": "What do my notes say about Docker?"
`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.Ask(
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
		"invalid request",
	)

	assert.False(
		t,
		service.Called,
	)

}

func TestChatHandlerAsk_EmptyMessage(t *testing.T) {
	service := &mocks.FakeChatService{}

	handler := NewChatHandler(service)

	body := `{
	"session_id": 1,
	"message": ""
}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.Ask(
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
		"message is required",
	)

	assert.False(
		t,
		service.Called,
	)
}

func TestChatHandlerAsk_ServiceError(t *testing.T) {
	service := &mocks.FakeChatService{
		Err: errors.New("openai unavailable"),
	}

	handler := NewChatHandler(service)

	body := `{
	"session_id": 1,
	"message": "What do my notes say about Docker?"
}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.Ask(
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
		"failed to process chat",
	)

	assert.True(
		t,
		service.Called,
	)

	assert.Equal(
		t,
		1,
		service.LastSessionID,
	)

	assert.Equal(
		t,
		"What do my notes say about Docker?",
		service.LastMessage,
	)

}

// TestChatHandlerAsk_ServiceErrorDoesNotLeakDetail proves that whatever the
// service layer returns (which can wrap raw OpenAI/DB error text) never
// reaches the client - only the generic message does, while the detail is
// still logged server-side.
func TestChatHandlerAsk_ServiceErrorDoesNotLeakDetail(t *testing.T) {
	sensitive := "sk-fake123: invalid api key"
	service := &mocks.FakeChatService{
		Err: errors.New(sensitive),
	}

	handler := NewChatHandler(service)

	body := `{
	"session_id": 1,
	"message": "Hello"
}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.Ask(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusInternalServerError,
		recorder.Code,
	)

	assert.NotContains(
		t,
		recorder.Body.String(),
		sensitive,
	)
}

func TestChatHandlerAsk_SessionNotFound(t *testing.T) {
	service := &mocks.FakeChatService{
		Err: repository.ErrChatSessionNotFound,
	}

	handler := NewChatHandler(service)

	body := `{
	"session_id": 999,
	"message": "Hello"
}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.Ask(
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

func TestChatHandlerAsk_Unauthorized(t *testing.T) {
	service := &mocks.FakeChatService{}

	handler := NewChatHandler(service)

	body := `{
	"session_id": 1,
	"message": "Hello"
}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.Ask(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusUnauthorized,
		recorder.Code,
	)

	assert.False(
		t,
		service.Called,
	)
}

func TestChatHandlerAsk_NoteSession(t *testing.T) {
	service := &mocks.FakeChatService{
		Answer:          "Go is a compiled language.",
		ResponseSession: 10,
	}

	handler := NewChatHandler(service)

	body := `{
	"session_id": 10,
	"message": "Explain this note"
}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)
	req = withUserID(req, testUserID)

	recorder := httptest.NewRecorder()

	handler.Ask(
		recorder,
		req,
	)

	require.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	var response models.ChatResponse

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&response)

	require.NoError(t, err)

	assert.Equal(
		t,
		"Go is a compiled language.",
		response.Answer,
	)

	assert.True(
		t,
		service.Called,
	)

	assert.Equal(
		t,
		10,
		service.LastSessionID,
	)

	assert.Equal(
		t,
		"Explain this note",
		service.LastMessage,
	)

}
