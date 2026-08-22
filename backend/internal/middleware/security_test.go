package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSecurityHeaders(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	recorder := httptest.NewRecorder()

	SecurityHeaders(next).ServeHTTP(recorder, req)

	assert.Equal(t, "nosniff", recorder.Header().Get("X-Content-Type-Options"))
	assert.Equal(t, "DENY", recorder.Header().Get("X-Frame-Options"))
	assert.Equal(t, "no-referrer", recorder.Header().Get("Referrer-Policy"))
	assert.Equal(t, "default-src 'none'", recorder.Header().Get("Content-Security-Policy"))
	assert.Equal(t, "no-store", recorder.Header().Get("Cache-Control"))
}
