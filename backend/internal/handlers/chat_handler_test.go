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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatHandlerAsk(t *testing.T) {

	service := &mocks.FakeChatService{
		Answer: "Docker is a container platform.",
	}

	handler := NewChatHandler(service)

	body := `{
		"message":"What do my notes say about Docker?"
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

	require.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	var response models.ChatResponse

	err := json.NewDecoder(recorder.Body).Decode(&response)

	require.NoError(t, err)

	assert.Equal(
		t,
		"Docker is a container platform.",
		response.Answer,
	)

	assert.Equal(
		t,
		"What do my notes say about Docker?",
		service.LastMessage,
	)
}

func TestChatHandlerAsk_InvalidJSON(t *testing.T) {

	service := &mocks.FakeChatService{}

	handler := NewChatHandler(service)

	body := `{
		"message":"What do my notes say about Docker?"
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

func TestChatHandlerAsk_ServiceError(t *testing.T) {

	service := &mocks.FakeChatService{
		Err: errors.New("openai unavailable"),
	}

	handler := NewChatHandler(service)

	body := `{
		"message":"What do my notes say about Docker?"
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
		"What do my notes say about Docker?",
		service.LastMessage,
	)
}
