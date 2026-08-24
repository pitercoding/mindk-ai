package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/testutil"
	"github.com/stretchr/testify/require"
)

// TestNotes_FullCRUDLifecycle drives a note through every stage of its life
// through the real HTTP stack - create, list, read, update, delete, then
// confirm it is truly gone - proving the full chain (middleware -> handler
// -> service -> repository -> SQLite), not just each layer in isolation.
func TestNotes_FullCRUDLifecycle(t *testing.T) {
	server := testutil.NewServer(t)
	token := testutil.TokenFor(t, "user_a")

	authed := func(method, url string, body []byte) *http.Response {
		var reader *bytes.Reader
		if body != nil {
			reader = bytes.NewReader(body)
		} else {
			reader = bytes.NewReader(nil)
		}

		req, err := http.NewRequest(method, url, reader)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)
		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		res, err := server.Client().Do(req)
		require.NoError(t, err)
		t.Cleanup(func() { res.Body.Close() })

		return res
	}

	// Create
	res := authed(http.MethodPost, server.URL+"/notes", []byte(`{"title":"Test note","content":"Hello MindK"}`))
	require.Equal(t, http.StatusCreated, res.StatusCode)

	var created models.Note
	require.NoError(t, json.NewDecoder(res.Body).Decode(&created))
	require.Equal(t, "Test note", created.Title)
	require.NotZero(t, created.ID)

	noteURL := fmt.Sprintf("%s/notes/%d", server.URL, created.ID)

	// List
	res = authed(http.MethodGet, server.URL+"/notes", nil)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var notes []models.Note
	require.NoError(t, json.NewDecoder(res.Body).Decode(&notes))
	require.Len(t, notes, 1)
	require.Equal(t, "Test note", notes[0].Title)

	// Read by ID
	res = authed(http.MethodGet, noteURL, nil)
	require.Equal(t, http.StatusOK, res.StatusCode)

	var fetched models.Note
	require.NoError(t, json.NewDecoder(res.Body).Decode(&fetched))
	require.Equal(t, created.ID, fetched.ID)

	// Update
	res = authed(http.MethodPut, noteURL, []byte(`{"title":"Updated title","content":"Updated content"}`))
	require.Equal(t, http.StatusOK, res.StatusCode)

	var updated models.Note
	require.NoError(t, json.NewDecoder(res.Body).Decode(&updated))
	require.Equal(t, "Updated title", updated.Title)

	// Confirm the update persisted
	res = authed(http.MethodGet, noteURL, nil)
	require.Equal(t, http.StatusOK, res.StatusCode)
	require.NoError(t, json.NewDecoder(res.Body).Decode(&fetched))
	require.Equal(t, "Updated title", fetched.Title)

	// Delete
	res = authed(http.MethodDelete, noteURL, nil)
	require.Equal(t, http.StatusNoContent, res.StatusCode)

	// Confirm it is truly gone
	res = authed(http.MethodGet, noteURL, nil)
	require.Equal(t, http.StatusNotFound, res.StatusCode)
}
