package httputil

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetIDFromPath(t *testing.T) {
	t.Run("returns the id at the end of the path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/notes/42", nil)

		id, err := GetIDFromPath(req)

		require.NoError(t, err)
		assert.Equal(t, 42, id)
	})

	t.Run("works for nested paths", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/chat/sessions/7", nil)

		id, err := GetIDFromPath(req)

		require.NoError(t, err)
		assert.Equal(t, 7, id)
	})

	t.Run("rejects invalid ids", func(t *testing.T) {
		for _, path := range []string{
			"/notes/abc",
			"/notes/1.5",
			"/notes/",
			"/notes/0",
			"/notes/-3",
		} {
			req := httptest.NewRequest(http.MethodGet, path, nil)

			id, err := GetIDFromPath(req)

			assert.Error(t, err, path)
			assert.Zero(t, id, path)
		}
	})
}
