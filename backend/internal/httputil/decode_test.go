package httputil

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type payload struct {
	Name string `json:"name"`
}

func TestDecodeJSON(t *testing.T) {
	t.Run("decodes a valid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"mindk"}`))
		rec := httptest.NewRecorder()

		var got payload
		err := DecodeJSON(rec, req, MaxJSONBodyBytes, &got)

		require.NoError(t, err)
		assert.Equal(t, "mindk", got.Name)
		assert.Equal(t, http.StatusOK, rec.Code, "nothing should be written on success")
	})

	t.Run("responds 400 for malformed JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{not json`))
		rec := httptest.NewRecorder()

		var got payload
		err := DecodeJSON(rec, req, MaxJSONBodyBytes, &got)

		require.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "invalid request body")
	})

	t.Run("responds 400 for an empty body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		rec := httptest.NewRecorder()

		var got payload
		err := DecodeJSON(rec, req, MaxJSONBodyBytes, &got)

		require.Error(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("responds 413 when the body exceeds the limit", func(t *testing.T) {
		body := `{"name":"` + strings.Repeat("a", 100) + `"}`
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		rec := httptest.NewRecorder()

		var got payload
		err := DecodeJSON(rec, req, 10, &got)

		require.Error(t, err)
		assert.Equal(t, http.StatusRequestEntityTooLarge, rec.Code)
		assert.Contains(t, rec.Body.String(), "request body too large")
	})
}
