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

// TestPostgres_FullStack proves the application works end to end against a
// real PostgreSQL database, not just SQLite: migrations apply cleanly to an
// empty schema, the $N placeholders and RETURNING id used by the
// repositories are valid PostgreSQL, and ownership filtering still holds.
// It is skipped unless POSTGRES_TEST_URL is set - see
// testutil.NewPostgresServer.
func TestPostgres_FullStack(t *testing.T) {
	server := testutil.NewPostgresServer(t)

	t.Run("health", func(t *testing.T) {
		res, err := server.Client().Get(server.URL + "/health")
		require.NoError(t, err)
		defer res.Body.Close()
		require.Equal(t, http.StatusOK, res.StatusCode)
	})

	tokens := testutil.TokensFor(t, "user_a", "user_b")
	tokenA, tokenB := tokens["user_a"], tokens["user_b"]

	t.Run("notes CRUD", func(t *testing.T) {
		res := authedRequest(t, server.Client(), tokenA, http.MethodPost, server.URL+"/notes",
			[]byte(`{"title":"pg note","content":"hello postgres"}`))
		require.Equal(t, http.StatusCreated, res.StatusCode)

		var created models.Note
		require.NoError(t, json.NewDecoder(res.Body).Decode(&created))
		require.NotZero(t, created.ID, "INSERT ... RETURNING id must populate the note ID")

		noteURL := fmt.Sprintf("%s/notes/%d", server.URL, created.ID)

		res = authedRequest(t, server.Client(), tokenA, http.MethodGet, server.URL+"/notes", nil)
		require.Equal(t, http.StatusOK, res.StatusCode)
		var notes []models.Note
		require.NoError(t, json.NewDecoder(res.Body).Decode(&notes))
		require.Len(t, notes, 1)

		res = authedRequest(t, server.Client(), tokenA, http.MethodPut, noteURL,
			[]byte(`{"title":"pg note updated","content":"still here"}`))
		require.Equal(t, http.StatusOK, res.StatusCode)
		var updated models.Note
		require.NoError(t, json.NewDecoder(res.Body).Decode(&updated))
		require.Equal(t, "pg note updated", updated.Title)

		res = authedRequest(t, server.Client(), tokenA, http.MethodDelete, noteURL, nil)
		require.Equal(t, http.StatusNoContent, res.StatusCode)

		res = authedRequest(t, server.Client(), tokenA, http.MethodGet, noteURL, nil)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	t.Run("documents CRUD and search", func(t *testing.T) {
		res := authedRequest(t, server.Client(), tokenA, http.MethodPost, server.URL+"/documents",
			[]byte(`{"name":"pg-doc.txt","type":".txt","content":"searchable postgres content"}`))
		require.Equal(t, http.StatusCreated, res.StatusCode)

		var created models.Document
		require.NoError(t, json.NewDecoder(res.Body).Decode(&created))
		require.NotZero(t, created.ID)

		docURL := fmt.Sprintf("%s/documents/%d", server.URL, created.ID)

		res = authedRequest(t, server.Client(), tokenA, http.MethodGet, server.URL+"/documents", nil)
		require.Equal(t, http.StatusOK, res.StatusCode)
		var documents []models.Document
		require.NoError(t, json.NewDecoder(res.Body).Decode(&documents))
		require.Len(t, documents, 1)

		res = authedRequest(t, server.Client(), tokenA, http.MethodGet,
			server.URL+"/documents/search?q=searchable", nil)
		require.Equal(t, http.StatusOK, res.StatusCode)
		var results []models.Document
		require.NoError(t, json.NewDecoder(res.Body).Decode(&results))
		require.Len(t, results, 1, "LIKE-based search must match PostgreSQL the same way it does SQLite")

		res = authedRequest(t, server.Client(), tokenA, http.MethodDelete, docURL, nil)
		require.Equal(t, http.StatusNoContent, res.StatusCode)

		res = authedRequest(t, server.Client(), tokenA, http.MethodGet, docURL, nil)
		require.Equal(t, http.StatusNotFound, res.StatusCode)
	})

	var sessionID int
	t.Run("chat sessions and messages CRUD", func(t *testing.T) {
		res := authedRequest(t, server.Client(), tokenA, http.MethodPost, server.URL+"/chat/sessions",
			[]byte(`{"title":"pg chat","mode":"knowledge"}`))
		require.Equal(t, http.StatusCreated, res.StatusCode)

		var session models.ChatSession
		require.NoError(t, json.NewDecoder(res.Body).Decode(&session))
		require.NotZero(t, session.ID)
		sessionID = session.ID

		res = authedRequest(t, server.Client(), tokenA, http.MethodGet, server.URL+"/chat/sessions", nil)
		require.Equal(t, http.StatusOK, res.StatusCode)
		var sessions []models.ChatSession
		require.NoError(t, json.NewDecoder(res.Body).Decode(&sessions))
		require.Len(t, sessions, 1)

		messagesURL := fmt.Sprintf("%s/chat/messages/%d", server.URL, sessionID)
		res = authedRequest(t, server.Client(), tokenA, http.MethodPost, messagesURL,
			[]byte(`{"session_id":`+fmt.Sprint(sessionID)+`,"role":"user","content":"hello from postgres"}`))
		require.Equal(t, http.StatusCreated, res.StatusCode)

		res = authedRequest(t, server.Client(), tokenA, http.MethodGet, messagesURL, nil)
		require.Equal(t, http.StatusOK, res.StatusCode)
		var messages []models.ChatMessage
		require.NoError(t, json.NewDecoder(res.Body).Decode(&messages))
		require.Len(t, messages, 1)
		require.Equal(t, "hello from postgres", messages[0].Content)

		sessionURL := fmt.Sprintf("%s/chat/sessions/%d", server.URL, sessionID)
		res = authedRequest(t, server.Client(), tokenA, http.MethodDelete, sessionURL, nil)
		require.Equal(t, http.StatusNoContent, res.StatusCode)
	})

	t.Run("ownership: user B cannot reach user A's notes", func(t *testing.T) {
		res := authedRequest(t, server.Client(), tokenA, http.MethodPost, server.URL+"/notes",
			[]byte(`{"title":"A's secret","content":"only A can see this"}`))
		require.Equal(t, http.StatusCreated, res.StatusCode)

		var note models.Note
		require.NoError(t, json.NewDecoder(res.Body).Decode(&note))
		noteURL := fmt.Sprintf("%s/notes/%d", server.URL, note.ID)

		res = authedRequest(t, server.Client(), tokenB, http.MethodGet, noteURL, nil)
		require.Equal(t, http.StatusNotFound, res.StatusCode)

		res = authedRequest(t, server.Client(), tokenB, http.MethodPut, noteURL,
			[]byte(`{"title":"hacked","content":"hacked"}`))
		require.Equal(t, http.StatusNotFound, res.StatusCode)

		res = authedRequest(t, server.Client(), tokenB, http.MethodDelete, noteURL, nil)
		require.Equal(t, http.StatusNotFound, res.StatusCode)

		res = authedRequest(t, server.Client(), tokenB, http.MethodGet, server.URL+"/notes", nil)
		require.Equal(t, http.StatusOK, res.StatusCode)
		var notesB []models.Note
		require.NoError(t, json.NewDecoder(res.Body).Decode(&notesB))
		require.Empty(t, notesB, "user B's list must not include user A's note")
	})
}
