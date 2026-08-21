package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pitercoding/mindk-ai/backend/internal/auth"
	"github.com/stretchr/testify/require"
)

func newRateLimitedTestServer(rl *RateLimiter, userID string) *httptest.Server {
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := auth.WithUserID(r.Context(), userID)
		handler.ServeHTTP(w, r.WithContext(ctx))
	}))

	return httptest.NewServer(mux)
}

func doRequest(t *testing.T, url string) *http.Response {
	t.Helper()

	res, err := http.Get(url)
	require.NoError(t, err)
	t.Cleanup(func() { res.Body.Close() })

	return res
}

func TestRateLimiter_AllowsRequestsWithinLimit(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)
	server := newRateLimitedTestServer(rl, "user_a")
	defer server.Close()

	for i := 0; i < 3; i++ {
		res := doRequest(t, server.URL)
		require.Equal(t, http.StatusOK, res.StatusCode)
	}
}

func TestRateLimiter_RejectsRequestsExceedingLimit(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)
	server := newRateLimitedTestServer(rl, "user_a")
	defer server.Close()

	for i := 0; i < 3; i++ {
		res := doRequest(t, server.URL)
		require.Equal(t, http.StatusOK, res.StatusCode)
	}

	res := doRequest(t, server.URL)
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	require.NotEmpty(t, res.Header.Get("Retry-After"))
}

func TestRateLimiter_TracksUsersIndependently(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)

	serverA := newRateLimitedTestServer(rl, "user_a")
	defer serverA.Close()
	serverB := newRateLimitedTestServer(rl, "user_b")
	defer serverB.Close()

	for i := 0; i < 3; i++ {
		res := doRequest(t, serverA.URL)
		require.Equal(t, http.StatusOK, res.StatusCode)
	}

	// user_a is now at their limit, but user_b has made no requests yet.
	res := doRequest(t, serverA.URL)
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)

	for i := 0; i < 3; i++ {
		res := doRequest(t, serverB.URL)
		require.Equal(t, http.StatusOK, res.StatusCode)
	}
}

func TestRateLimiter_ResetsAfterWindowElapses(t *testing.T) {
	rl := NewRateLimiter(2, time.Minute)

	current := time.Now()
	rl.now = func() time.Time { return current }

	server := newRateLimitedTestServer(rl, "user_a")
	defer server.Close()

	require.Equal(t, http.StatusOK, doRequest(t, server.URL).StatusCode)
	require.Equal(t, http.StatusOK, doRequest(t, server.URL).StatusCode)
	require.Equal(t, http.StatusTooManyRequests, doRequest(t, server.URL).StatusCode)

	// Advance the injected clock past the window instead of sleeping.
	current = current.Add(time.Minute + time.Second)

	require.Equal(t, http.StatusOK, doRequest(t, server.URL).StatusCode)
}

func TestRateLimiter_RejectsUnauthenticatedRequests(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)
	handler := rl.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	require.Equal(t, http.StatusUnauthorized, w.Code)
}
