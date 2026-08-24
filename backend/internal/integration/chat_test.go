package integration

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/routes"
	"github.com/pitercoding/mindk-ai/backend/internal/testutil"
	"github.com/stretchr/testify/require"
)

// TestChatSessions_FullCRUDLifecycle drives a chat session through create,
// list, read, update and delete through the real HTTP stack.
func TestChatSessions_FullCRUDLifecycle(t *testing.T) {
	server := testutil.NewServer(t)
	token := testutil.TokenFor(t, "user_a")

	res := authedRequest(t, server.Client(), token, http.MethodPost, server.URL+"/chat/sessions",
		[]byte(`{"title":"My chat","mode":"knowledge"}`))
	require.Equal(t, http.StatusCreated, res.StatusCode)

	var created models.ChatSession
	require.NoError(t, json.NewDecoder(res.Body).Decode(&created))
	require.Equal(t, "My chat", created.Title)

	sessionURL := fmt.Sprintf("%s/chat/sessions/%d", server.URL, created.ID)

	res = authedRequest(t, server.Client(), token, http.MethodGet, server.URL+"/chat/sessions", nil)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var sessions []models.ChatSession
	require.NoError(t, json.NewDecoder(res.Body).Decode(&sessions))
	require.Len(t, sessions, 1)

	res = authedRequest(t, server.Client(), token, http.MethodPut, sessionURL,
		[]byte(`{"title":"Renamed chat","mode":"knowledge"}`))
	require.Equal(t, http.StatusOK, res.StatusCode)

	var updated models.ChatSession
	require.NoError(t, json.NewDecoder(res.Body).Decode(&updated))
	require.Equal(t, "Renamed chat", updated.Title)

	res = authedRequest(t, server.Client(), token, http.MethodDelete, sessionURL, nil)
	require.Equal(t, http.StatusNoContent, res.StatusCode)

	res = authedRequest(t, server.Client(), token, http.MethodGet, sessionURL, nil)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

// TestChatMessages_SaveAndRetrieveWithinSession proves a message posted to a
// session is retrievable afterwards, exercising the rule that a message
// belongs to a session which in turn belongs to a user.
func TestChatMessages_SaveAndRetrieveWithinSession(t *testing.T) {
	server := testutil.NewServer(t)
	token := testutil.TokenFor(t, "user_a")

	res := authedRequest(t, server.Client(), token, http.MethodPost, server.URL+"/chat/sessions",
		[]byte(`{"title":"My chat","mode":"knowledge"}`))
	require.Equal(t, http.StatusCreated, res.StatusCode)

	var session models.ChatSession
	require.NoError(t, json.NewDecoder(res.Body).Decode(&session))

	messagesURL := fmt.Sprintf("%s/chat/messages/%d", server.URL, session.ID)

	res = authedRequest(t, server.Client(), token, http.MethodPost, messagesURL,
		[]byte(`{"session_id":`+fmt.Sprint(session.ID)+`,"role":"user","content":"hello there"}`))
	require.Equal(t, http.StatusCreated, res.StatusCode)

	res = authedRequest(t, server.Client(), token, http.MethodGet, messagesURL, nil)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var messages []models.ChatMessage
	require.NoError(t, json.NewDecoder(res.Body).Decode(&messages))
	require.Len(t, messages, 1)
	require.Equal(t, "hello there", messages[0].Content)
}

// TestChat_AskUsesFakeLLMAndPersistsHistory proves POST /chat runs through
// the full stack - auto-creating a session, calling the (fake) LLM and
// saving both the user's message and the assistant's answer - without ever
// making a real OpenAI call.
func TestChat_AskUsesFakeLLMAndPersistsHistory(t *testing.T) {
	server := testutil.NewServer(t)
	server.LLM.Response = "42 is the answer"
	token := testutil.TokenFor(t, "user_a")

	res := authedRequest(t, server.Client(), token, http.MethodPost, server.URL+"/chat",
		[]byte(`{"message":"what is the answer?","mode":"knowledge"}`))
	require.Equal(t, http.StatusOK, res.StatusCode)

	var chatRes models.ChatResponse
	require.NoError(t, json.NewDecoder(res.Body).Decode(&chatRes))
	require.Equal(t, "42 is the answer", chatRes.Answer)
	require.NotZero(t, chatRes.SessionID)

	require.Contains(t, server.LLM.LastPrompt, "what is the answer?")

	messagesURL := fmt.Sprintf("%s/chat/messages/%d", server.URL, chatRes.SessionID)
	res = authedRequest(t, server.Client(), token, http.MethodGet, messagesURL, nil)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var messages []models.ChatMessage
	require.NoError(t, json.NewDecoder(res.Body).Decode(&messages))
	require.Len(t, messages, 2, "both the user's message and the assistant's answer must be persisted")
	require.Equal(t, "user", messages[0].Role)
	require.Equal(t, "what is the answer?", messages[0].Content)
	require.Equal(t, "assistant", messages[1].Role)
	require.Equal(t, "42 is the answer", messages[1].Content)
}

// TestChat_RateLimitIsStricterThanGeneralAPI proves /chat is throttled at
// routes.ChatRateLimit, distinct from (and lower than) the general API
// limit that governs notes/documents/sessions.
func TestChat_RateLimitIsStricterThanGeneralAPI(t *testing.T) {
	server := testutil.NewServer(t)
	token := testutil.TokenFor(t, "user_a")

	ask := func() *http.Response {
		return authedRequest(t, server.Client(), token, http.MethodPost, server.URL+"/chat",
			[]byte(`{"message":"hi","mode":"knowledge"}`))
	}

	for i := 0; i < routes.ChatRateLimit; i++ {
		res := ask()
		require.Equal(t, http.StatusOK, res.StatusCode, "request %d should be within the chat rate limit", i+1)
	}

	res := ask()
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	require.NotEmpty(t, res.Header.Get("Retry-After"))
}

// TestChat_LLMFailureReturnsGenericError proves that when the LLM call
// fails, the client sees a generic 500 message - not the underlying error
// detail, which could otherwise leak internal infrastructure information.
func TestChat_LLMFailureReturnsGenericError(t *testing.T) {
	server := testutil.NewServer(t)
	server.LLM.Err = errors.New("dial tcp 203.0.113.5:443: connection refused (api.openai.com)")
	token := testutil.TokenFor(t, "user_a")

	res := authedRequest(t, server.Client(), token, http.MethodPost, server.URL+"/chat",
		[]byte(`{"message":"hi","mode":"knowledge"}`))
	require.Equal(t, http.StatusInternalServerError, res.StatusCode)

	body := readBody(t, res)
	require.Contains(t, body, "failed to process chat request")
	require.NotContains(t, body, "openai")
	require.NotContains(t, body, "dial tcp")
	require.NotContains(t, body, "203.0.113.5")
}
