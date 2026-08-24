package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureSlog swaps slog's default logger for one writing JSON lines into
// buf, restoring the previous default when the test ends.
func captureSlog(t *testing.T) *bytes.Buffer {
	t.Helper()

	buf := &bytes.Buffer{}
	previous := slog.Default()
	slog.SetDefault(slog.New(slog.NewJSONHandler(buf, nil)))
	t.Cleanup(func() { slog.SetDefault(previous) })

	return buf
}

// decodeLastLogLine parses the last JSON log line written to buf.
func decodeLastLogLine(t *testing.T, buf *bytes.Buffer) map[string]any {
	t.Helper()

	lines := bytes.Split(bytes.TrimSpace(buf.Bytes()), []byte("\n"))
	require.NotEmpty(t, lines)

	var entry map[string]any
	require.NoError(t, json.Unmarshal(lines[len(lines)-1], &entry))

	return entry
}

func TestRequestLogger_SetsRequestIDHeader(t *testing.T) {
	captureSlog(t)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	RequestLogger(next).ServeHTTP(recorder, req)

	assert.NotEmpty(t, recorder.Header().Get("X-Request-ID"))
}

func TestRequestLogger_LogsMethodPathStatusAndDuration(t *testing.T) {
	buf := captureSlog(t)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	req := httptest.NewRequest(http.MethodGet, "/notes/999", nil)
	recorder := httptest.NewRecorder()

	RequestLogger(next).ServeHTTP(recorder, req)

	entry := decodeLastLogLine(t, buf)

	assert.Equal(t, "request completed", entry["msg"])
	assert.Equal(t, http.MethodGet, entry["method"])
	assert.Equal(t, "/notes/999", entry["path"])
	assert.Equal(t, float64(http.StatusNotFound), entry["status"])
	assert.Equal(t, recorder.Header().Get("X-Request-ID"), entry["request_id"])
	assert.Contains(t, entry, "duration_ms")
	assert.NotContains(t, entry, "user_id", "unauthenticated requests must not log a user_id")
}

func TestRequestLogger_DefaultsToStatus200WhenHandlerNeverCallsWriteHeader(t *testing.T) {
	buf := captureSlog(t)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	RequestLogger(next).ServeHTTP(recorder, req)

	entry := decodeLastLogLine(t, buf)
	assert.Equal(t, float64(http.StatusOK), entry["status"])
}

func TestRequestLogger_IncludesUserIDWhenClerkRecordsIt(t *testing.T) {
	buf := captureSlog(t)

	// Mirrors what Clerk auth does: record the user ID for the access log,
	// then also put it in context the normal way for the handler itself.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recordUserID(r.Context(), "user_abc")
		ctx := auth.WithUserID(r.Context(), "user_abc")
		r = r.WithContext(ctx)
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	recorder := httptest.NewRecorder()

	RequestLogger(next).ServeHTTP(recorder, req)

	entry := decodeLastLogLine(t, buf)
	assert.Equal(t, "user_abc", entry["user_id"])
}

// TestRequestLogger_NeverLogsAuthorizationHeader proves the access log line
// carries no trace of the raw Authorization header, even though it is
// present on the request RequestLogger wraps.
func TestRequestLogger_NeverLogsAuthorizationHeader(t *testing.T) {
	buf := captureSlog(t)

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	})

	req := httptest.NewRequest(http.MethodGet, "/notes", nil)
	req.Header.Set("Authorization", "Bearer super-secret-token")
	recorder := httptest.NewRecorder()

	RequestLogger(next).ServeHTTP(recorder, req)

	assert.NotContains(t, buf.String(), "super-secret-token")
	assert.NotContains(t, buf.String(), "Bearer")
}
