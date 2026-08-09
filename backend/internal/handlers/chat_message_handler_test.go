package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/services/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatMessageHandlerSave(t *testing.T) {
	service := &mocks.FakeChatMessageService{}

	handler := NewChatMessageHandler(service)

	body := `{
	"session_id": 1,
	"role": "user",
	"content": "What is Go?"
}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat/messages/1",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.Save(
		recorder,
		req,
	)

	require.Equal(
		t,
		http.StatusCreated,
		recorder.Code,
	)

	require.Len(
		t,
		service.Saved,
		1,
	)

	assert.Equal(
		t,
		1,
		service.Saved[0].SessionID,
	)

	assert.Equal(
		t,
		"user",
		service.Saved[0].Role,
	)

	assert.Equal(
		t,
		"What is Go?",
		service.Saved[0].Content,
	)

	var response models.ChatMessage

	err := json.NewDecoder(
		recorder.Body,
	).Decode(&response)

	require.NoError(t, err)

	assert.Equal(
		t,
		1,
		response.SessionID,
	)

	assert.Equal(
		t,
		"user",
		response.Role,
	)

	assert.Equal(
		t,
		"What is Go?",
		response.Content,
	)

}

func TestChatMessageHandlerSave_InvalidJSON(t *testing.T) {
	service := &mocks.FakeChatMessageService{}

	handler := NewChatMessageHandler(service)

	body := `{
	"session_id": 1,
	"role": "user",
	"content": "What is Go?"
`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat/messages/1",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.Save(
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

	assert.Len(
		t,
		service.Saved,
		0,
	)

}

func TestChatMessageHandlerSave_EmptyContent(t *testing.T) {
	service := &mocks.FakeChatMessageService{}

	handler := NewChatMessageHandler(service)

	body := `{
	"session_id": 1,
	"role": "user",
	"content": ""
}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat/messages/1",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.Save(
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
		"content is required",
	)

	assert.Len(
		t,
		service.Saved,
		0,
	)

}

func TestChatMessageHandlerSave_ServiceError(t *testing.T) {
	service := &mocks.FakeChatMessageService{
		Err: errors.New("database error"),
	}

	handler := NewChatMessageHandler(service)

	body := `{
	"session_id": 1,
	"role": "user",
	"content": "What is Go?"
}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/chat/messages/1",
		strings.NewReader(body),
	)

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	recorder := httptest.NewRecorder()

	handler.Save(
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
		"failed to save message",
	)

	require.Len(
		t,
		service.Saved,
		1,
	)

}

func TestChatMessageHandlerGetBySessionID(t *testing.T) {
	service := &mocks.FakeChatMessageService{
		Messages: []models.ChatMessage{
			{
				ID:        1,
				SessionID: 10,
				Role:      "user",
				Content:   "What is Go?",
				CreatedAt: time.Now(),
			},
			{
				ID:        2,
				SessionID: 10,
				Role:      "assistant",
				Content:   "Go is a programming language.",
				CreatedAt: time.Now(),
			},
		},
	}

	handler := NewChatMessageHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/chat/messages/10",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetBySessionID(
		recorder,
		req,
	)

	require.Equal(
		t,
		http.StatusOK,
		recorder.Code,
	)

	assert.Equal(
		t,
		10,
		service.LastSessionID,
	)

	var response []models.ChatMessage

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
		1,
		response[0].ID,
	)

	assert.Equal(
		t,
		10,
		response[0].SessionID,
	)

	assert.Equal(
		t,
		"user",
		response[0].Role,
	)

	assert.Equal(
		t,
		"What is Go?",
		response[0].Content,
	)

	assert.Equal(
		t,
		2,
		response[1].ID,
	)

	assert.Equal(
		t,
		"assistant",
		response[1].Role,
	)

	assert.Equal(
		t,
		"Go is a programming language.",
		response[1].Content,
	)

}

func TestChatMessageHandlerGetBySessionID_InvalidID(t *testing.T) {
	service := &mocks.FakeChatMessageService{}

	handler := NewChatMessageHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/chat/messages/invalid",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetBySessionID(
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

	assert.Equal(
		t,
		0,
		service.LastSessionID,
	)

}

func TestChatMessageHandlerGetBySessionID_ServiceError(t *testing.T) {
	service := &mocks.FakeChatMessageService{
		Err: errors.New("database error"),
	}

	handler := NewChatMessageHandler(service)

	req := httptest.NewRequest(
		http.MethodGet,
		"/chat/messages/10",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.GetBySessionID(
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
		"failed to fetch chat messages",
	)

	assert.Equal(
		t,
		10,
		service.LastSessionID,
	)

}

func TestChatMessageHandlerHandleMessages_MethodNotAllowed(t *testing.T) {
	service := &mocks.FakeChatMessageService{}

	handler := NewChatMessageHandler(service)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/chat/messages/10",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.HandleMessages(
		recorder,
		req,
	)

	assert.Equal(
		t,
		http.StatusMethodNotAllowed,
		recorder.Code,
	)

	assert.Contains(
		t,
		recorder.Body.String(),
		"method not allowed",
	)

}
