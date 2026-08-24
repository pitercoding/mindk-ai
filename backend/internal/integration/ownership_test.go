package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/testutil"
	"github.com/stretchr/testify/require"
)

// TestOwnership_NotesAreIsolatedBetweenUsers proves user B cannot read,
// update or delete a note created by user A - each operation must resolve
// to 404, matching the repository's "WHERE id = $1 AND user_id = $2" filter,
// never a 200 or a 403 that would confirm the note's existence.
func TestOwnership_NotesAreIsolatedBetweenUsers(t *testing.T) {
	server := testutil.NewServer(t)
	tokens := testutil.TokensFor(t, "user_a", "user_b")

	res := authedRequest(t, server.Client(), tokens["user_a"], http.MethodPost, server.URL+"/notes",
		[]byte(`{"title":"A's note","content":"secret"}`))
	require.Equal(t, http.StatusCreated, res.StatusCode)

	var note models.Note
	require.NoError(t, json.NewDecoder(res.Body).Decode(&note))

	noteURL := fmt.Sprintf("%s/notes/%d", server.URL, note.ID)

	res = authedRequest(t, server.Client(), tokens["user_b"], http.MethodGet, noteURL, nil)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	res = authedRequest(t, server.Client(), tokens["user_b"], http.MethodPut, noteURL,
		[]byte(`{"title":"hijacked","content":"hijacked"}`))
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	res = authedRequest(t, server.Client(), tokens["user_b"], http.MethodDelete, noteURL, nil)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	// The note must be untouched: user A can still read it as originally written.
	res = authedRequest(t, server.Client(), tokens["user_a"], http.MethodGet, noteURL, nil)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NoError(t, json.NewDecoder(res.Body).Decode(&note))
	require.Equal(t, "A's note", note.Title)
}

// TestOwnership_DocumentsAreIsolatedBetweenUsers mirrors the notes ownership
// test for documents: user B gets 404 for read and delete of user A's
// document.
func TestOwnership_DocumentsAreIsolatedBetweenUsers(t *testing.T) {
	server := testutil.NewServer(t)
	tokens := testutil.TokensFor(t, "user_a", "user_b")

	res := authedRequest(t, server.Client(), tokens["user_a"], http.MethodPost, server.URL+"/documents",
		[]byte(`{"name":"a-doc.txt","type":".txt","content":"A's private content"}`))
	require.Equal(t, http.StatusCreated, res.StatusCode)

	var doc models.Document
	require.NoError(t, json.NewDecoder(res.Body).Decode(&doc))

	docURL := fmt.Sprintf("%s/documents/%d", server.URL, doc.ID)

	res = authedRequest(t, server.Client(), tokens["user_b"], http.MethodGet, docURL, nil)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	res = authedRequest(t, server.Client(), tokens["user_b"], http.MethodDelete, docURL, nil)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	res = authedRequest(t, server.Client(), tokens["user_b"], http.MethodGet, server.URL+"/documents", nil)
	require.Equal(t, http.StatusOK, res.StatusCode)
	var userBDocs []models.Document
	require.NoError(t, json.NewDecoder(res.Body).Decode(&userBDocs))
	require.Empty(t, userBDocs, "user B must not see user A's documents in their own list")
}

// TestOwnership_ChatSessionsAndMessagesAreIsolatedBetweenUsers proves a chat
// session (and its messages) created by user A is invisible to user B at
// every layer: the session itself, and the messages within it.
func TestOwnership_ChatSessionsAndMessagesAreIsolatedBetweenUsers(t *testing.T) {
	server := testutil.NewServer(t)
	tokens := testutil.TokensFor(t, "user_a", "user_b")

	res := authedRequest(t, server.Client(), tokens["user_a"], http.MethodPost, server.URL+"/chat/sessions",
		[]byte(`{"title":"A's session","mode":"knowledge"}`))
	require.Equal(t, http.StatusCreated, res.StatusCode)

	var session models.ChatSession
	require.NoError(t, json.NewDecoder(res.Body).Decode(&session))

	messagesURL := fmt.Sprintf("%s/chat/messages/%d", server.URL, session.ID)
	res = authedRequest(t, server.Client(), tokens["user_a"], http.MethodPost, messagesURL,
		[]byte(fmt.Sprintf(`{"session_id":%d,"role":"user","content":"A's secret question"}`, session.ID)))
	require.Equal(t, http.StatusCreated, res.StatusCode)

	sessionURL := fmt.Sprintf("%s/chat/sessions/%d", server.URL, session.ID)

	// User B cannot see the session itself.
	res = authedRequest(t, server.Client(), tokens["user_b"], http.MethodGet, sessionURL, nil)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	// Nor read its messages...
	res = authedRequest(t, server.Client(), tokens["user_b"], http.MethodGet, messagesURL, nil)
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	// ...nor post a message into it.
	res = authedRequest(t, server.Client(), tokens["user_b"], http.MethodPost, messagesURL,
		[]byte(fmt.Sprintf(`{"session_id":%d,"role":"user","content":"trying to inject"}`, session.ID)))
	require.Equal(t, http.StatusNotFound, res.StatusCode)

	// Nor delete it.
	res = authedRequest(t, server.Client(), tokens["user_b"], http.MethodDelete, sessionURL, nil)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}
