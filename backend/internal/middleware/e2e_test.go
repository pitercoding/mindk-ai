package middleware

import (
	"bytes"
	"crypto/rsa"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/pitercoding/mindk-ai/backend/internal/handlers"
	"github.com/pitercoding/mindk-ai/backend/internal/migrations"
	"github.com/pitercoding/mindk-ai/backend/internal/models"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
	"github.com/pitercoding/mindk-ai/backend/internal/services"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

// newE2ETestServer wires the real Clerk middleware in front of the real
// NoteHandler, NoteService and NoteRepository (backed by an in-memory
// SQLite DB), mirroring exactly how routes.go wires /notes in production.
// This proves the full chain: JWT -> Clerk -> claims.Subject -> context ->
// handler -> service -> repository -> SQL ownership filter.
func newE2ETestServer(t *testing.T) *httptest.Server {
	t.Helper()

	return newE2ETestServerWithRateLimit(t, nil)
}

// newE2ETestServerWithRateLimit is newE2ETestServer with an optional rate
// limiter spliced in between Clerk and the NoteHandler, mirroring how
// routes.go wires /notes in production: Clerk -> RateLimit -> handler.
// Pass a nil limiter to skip rate limiting entirely.
func newE2ETestServerWithRateLimit(t *testing.T, rl *RateLimiter) *httptest.Server {
	t.Helper()

	db, err := sql.Open("sqlite", "file::memory:?cache=shared")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	// A single connection keeps every query on the same in-memory database.
	db.SetMaxOpenConns(1)
	require.NoError(t, migrations.Run(db))

	noteRepo := repository.NewNoteRepository(db)
	noteService := services.NewNoteService(noteRepo)
	noteHandler := handlers.NewNoteHandler(noteService)

	authMiddleware := NewClerkAuth("test_secret_key")

	wrap := func(h http.HandlerFunc) http.Handler {
		var handler http.Handler = h
		if rl != nil {
			handler = rl.Middleware(handler)
		}
		return authMiddleware(handler)
	}

	mux := http.NewServeMux()
	mux.Handle("/notes", wrap(noteHandler.HandleNotes))
	mux.Handle("/notes/", wrap(noteHandler.HandleNote))

	server := httptest.NewServer(RequestLogger(SecurityHeaders(mux)))
	t.Cleanup(server.Close)

	return server
}

// generateToken creates a valid Clerk JWT for the given user, without
// yet exposing its public key over a JWKS endpoint. Use tokenFor for a
// single-user test, or startMockJWKSMulti to serve several users' keys
// from one JWKS endpoint when a test needs more than one live token at
// once (the Clerk backend is a single global pointer, so only the most
// recently started mock JWKS server is reachable for verification).
func generateToken(t *testing.T, userID string) (token, kid string, pub *rsa.PublicKey) {
	t.Helper()

	kid = "kid-" + t.Name() + "-" + userID
	tok, rawPub := clerktest.GenerateJWT(t, map[string]any{
		"sub": userID,
		"sid": "sess_" + userID,
		"iss": "https://clerk.test.dev",
	}, kid)

	rsaPub, ok := rawPub.(*rsa.PublicKey)
	require.True(t, ok)

	return tok, kid, rsaPub
}

// tokenFor generates a valid Clerk JWT for the given user and starts the
// mock JWKS endpoint needed to verify it, exactly like clerk_test.go does.
// Only use this when a test needs a single live token at a time.
func tokenFor(t *testing.T, userID string) string {
	t.Helper()

	token, kid, pub := generateToken(t, userID)
	startMockJWKS(t, kid, pub)

	return token
}

// startMockJWKSMulti serves the public keys of several users from one
// JWKS endpoint, so tokens for all of them can be verified concurrently
// against the single global Clerk backend pointer.
func startMockJWKSMulti(t *testing.T, keys map[string]*rsa.PublicKey) {
	t.Helper()

	type jwk struct {
		Use string `json:"use"`
		Kty string `json:"kty"`
		Kid string `json:"kid"`
		Alg string `json:"alg"`
		N   string `json:"n"`
		E   string `json:"e"`
	}

	var jwks struct {
		Keys []jwk `json:"keys"`
	}
	for kid, pub := range keys {
		jwks.Keys = append(jwks.Keys, jwk{
			Use: "sig",
			Kty: "RSA",
			Kid: kid,
			Alg: "RS256",
			N:   base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
			E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
		})
	}

	body, err := json.Marshal(jwks)
	require.NoError(t, err)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/jwks" && r.Method == http.MethodGet {
			_, err := w.Write(body)
			require.NoError(t, err)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	t.Cleanup(server.Close)

	clerk.SetBackend(clerk.NewBackend(&clerk.BackendConfig{
		HTTPClient: server.Client(),
		URL:        &server.URL,
	}))
}

// TestE2E_NotesRequireAuthentication proves that /notes cannot be reached
// without a valid Clerk token: the middleware rejects the request before
// the handler (and thus the DB) is ever touched.
func TestE2E_NotesRequireAuthentication(t *testing.T) {
	server := newE2ETestServer(t)

	req, err := http.NewRequest(http.MethodPost, server.URL+"/notes", bytes.NewBufferString(`{"title":"t","content":"c"}`))
	require.NoError(t, err)

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

// TestE2E_NoteOwnershipComesFromToken_NotBody proves that even when a
// caller tries to spoof ownership via the request body, the note is always
// created for the user identified by the JWT - never the body's user_id.
func TestE2E_NoteOwnershipComesFromToken_NotBody(t *testing.T) {
	server := newE2ETestServer(t)
	token := tokenFor(t, "user_a")

	payload := `{"title":"real title","content":"real content","user_id":"user_b_spoofed"}`
	req, err := http.NewRequest(http.MethodPost, server.URL+"/notes", bytes.NewBufferString(payload))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusCreated, res.StatusCode)

	var created models.Note
	require.NoError(t, json.NewDecoder(res.Body).Decode(&created))
	require.Equal(t, "user_a", created.UserID, "the note must belong to the JWT's subject, not the spoofed body field")
}

// TestE2E_NotesAreIsolatedByAuthenticatedUser proves cross-user isolation
// through the full stack, driven purely by each caller's own token: user A
// never sees notes created by user B and vice versa.
func TestE2E_NotesAreIsolatedByAuthenticatedUser(t *testing.T) {
	server := newE2ETestServer(t)

	tokenA, kidA, pubA := generateToken(t, "user_a")
	tokenB, kidB, pubB := generateToken(t, "user_b")
	startMockJWKSMulti(t, map[string]*rsa.PublicKey{kidA: pubA, kidB: pubB})

	createNote := func(token, title string) {
		payload := `{"title":"` + title + `","content":"content"}`
		req, err := http.NewRequest(http.MethodPost, server.URL+"/notes", bytes.NewBufferString(payload))
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		res, err := server.Client().Do(req)
		require.NoError(t, err)
		defer res.Body.Close()
		require.Equal(t, http.StatusCreated, res.StatusCode)
	}

	createNote(tokenA, "A's note")
	createNote(tokenB, "B's note")

	req, err := http.NewRequest(http.MethodGet, server.URL+"/notes", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+tokenA)

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)

	var notes []models.Note
	require.NoError(t, json.NewDecoder(res.Body).Decode(&notes))
	require.Len(t, notes, 1, "user A must only see their own notes")
	require.Equal(t, "A's note", notes[0].Title)
	require.Equal(t, "user_a", notes[0].UserID)
}

// TestE2E_RateLimitAppliesAfterAuthentication proves the rate limiter is
// keyed by the JWT's subject rather than anything the client controls: it
// runs the real Clerk -> RateLimit -> NoteHandler chain and shows the
// caller is cut off with 429 once they exceed the limit for their token.
func TestE2E_RateLimitAppliesAfterAuthentication(t *testing.T) {
	rl := NewRateLimiter(3, time.Minute)
	server := newE2ETestServerWithRateLimit(t, rl)
	token := tokenFor(t, "user_a")

	getNotes := func() *http.Response {
		req, err := http.NewRequest(http.MethodGet, server.URL+"/notes", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		res, err := server.Client().Do(req)
		require.NoError(t, err)
		t.Cleanup(func() { res.Body.Close() })

		return res
	}

	for i := 0; i < 3; i++ {
		require.Equal(t, http.StatusOK, getNotes().StatusCode)
	}

	res := getNotes()
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)
	require.NotEmpty(t, res.Header.Get("Retry-After"))
}

// TestE2E_RateLimitIsIndependentPerAuthenticatedUser proves the limiter's
// key comes from claims.Subject, not from the client: two different users
// hitting the same endpoint through the same limiter never affect each
// other's quota, even though both requests arrive at the same server.
func TestE2E_RateLimitIsIndependentPerAuthenticatedUser(t *testing.T) {
	rl := NewRateLimiter(2, time.Minute)
	server := newE2ETestServerWithRateLimit(t, rl)

	tokenA, kidA, pubA := generateToken(t, "user_a")
	tokenB, kidB, pubB := generateToken(t, "user_b")
	startMockJWKSMulti(t, map[string]*rsa.PublicKey{kidA: pubA, kidB: pubB})

	getNotes := func(token string) *http.Response {
		req, err := http.NewRequest(http.MethodGet, server.URL+"/notes", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+token)

		res, err := server.Client().Do(req)
		require.NoError(t, err)
		t.Cleanup(func() { res.Body.Close() })

		return res
	}

	require.Equal(t, http.StatusOK, getNotes(tokenA).StatusCode)
	require.Equal(t, http.StatusOK, getNotes(tokenA).StatusCode)
	require.Equal(t, http.StatusTooManyRequests, getNotes(tokenA).StatusCode, "user A must now be rate limited")

	require.Equal(t, http.StatusOK, getNotes(tokenB).StatusCode, "user B's quota must be untouched by user A's requests")
	require.Equal(t, http.StatusOK, getNotes(tokenB).StatusCode)
}

// TestE2E_SecurityHeadersPresent proves SecurityHeaders is actually wired in
// front of the real stack, not just unit-tested in isolation - even an
// unauthenticated 401 response carries the headers.
func TestE2E_SecurityHeadersPresent(t *testing.T) {
	server := newE2ETestServer(t)

	res, err := server.Client().Get(server.URL + "/notes")
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.Equal(t, "nosniff", res.Header.Get("X-Content-Type-Options"))
	require.Equal(t, "DENY", res.Header.Get("X-Frame-Options"))
	require.Equal(t, "no-referrer", res.Header.Get("Referrer-Policy"))
	require.Equal(t, "default-src 'none'", res.Header.Get("Content-Security-Policy"))
	require.Equal(t, "no-store", res.Header.Get("Cache-Control"))
}

// TestE2E_OversizedBodyRejected proves a body bigger than the JSON limit is
// rejected with 413 before ever reaching the service/DB layer, driven
// through the real Clerk -> NoteHandler chain.
func TestE2E_OversizedBodyRejected(t *testing.T) {
	server := newE2ETestServer(t)
	token := tokenFor(t, "user_a")

	oversized := bytes.Repeat([]byte("a"), 2<<20) // 2 MiB, over the 1 MiB JSON limit
	payload := []byte(`{"title":"t","content":"`)
	payload = append(payload, oversized...)
	payload = append(payload, []byte(`"}`)...)

	req, err := http.NewRequest(http.MethodPost, server.URL+"/notes", bytes.NewReader(payload))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusRequestEntityTooLarge, res.StatusCode)
}

// TestE2E_MalformedJSONRejected proves invalid JSON never reaches the
// service layer and is rejected with 400.
func TestE2E_MalformedJSONRejected(t *testing.T) {
	server := newE2ETestServer(t)
	token := tokenFor(t, "user_a")

	req, err := http.NewRequest(http.MethodPost, server.URL+"/notes", bytes.NewBufferString(`{"title":"t"`))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
}

// TestE2E_InvalidInputRejected proves a request with an empty required
// field is rejected with 400 rather than reaching the DB.
func TestE2E_InvalidInputRejected(t *testing.T) {
	server := newE2ETestServer(t)
	token := tokenFor(t, "user_a")

	req, err := http.NewRequest(http.MethodPost, server.URL+"/notes", bytes.NewBufferString(`{"title":"","content":"c"}`))
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode)
}

// TestE2E_NonexistentResourceReturnsNotFound proves that a well-formed
// request for a note that doesn't exist resolves to 404, never a 500 that
// might otherwise leak internal error detail.
func TestE2E_NonexistentResourceReturnsNotFound(t *testing.T) {
	server := newE2ETestServer(t)
	token := tokenFor(t, "user_a")

	req, err := http.NewRequest(http.MethodGet, server.URL+"/notes/999", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusNotFound, res.StatusCode)
}

// TestE2E_RequestIDHeaderPresent proves RequestLogger is actually wired in
// front of the real stack: even an unauthenticated 401 response carries a
// correlation ID, so a client-reported error can always be traced server-side.
func TestE2E_RequestIDHeaderPresent(t *testing.T) {
	server := newE2ETestServer(t)

	res, err := server.Client().Get(server.URL + "/notes")
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	require.NotEmpty(t, res.Header.Get("X-Request-ID"))
}

// TestE2E_RequestLogIncludesUserIDAndMatchesResponseHeader proves the
// access log for an authenticated request carries the same request_id
// returned to the client, plus the user ID resolved from the JWT - not
// anything client-supplied - and nothing else about the request/response
// bodies.
func TestE2E_RequestLogIncludesUserIDAndMatchesResponseHeader(t *testing.T) {
	buf := captureSlog(t)

	server := newE2ETestServer(t)
	token := tokenFor(t, "user_a")

	req, err := http.NewRequest(http.MethodGet, server.URL+"/notes", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusOK, res.StatusCode)

	entry := decodeLastLogLine(t, buf)
	require.Equal(t, "request completed", entry["msg"])
	require.Equal(t, res.Header.Get("X-Request-ID"), entry["request_id"])
	require.Equal(t, "user_a", entry["user_id"])
	require.Equal(t, "/notes", entry["path"])
	require.Equal(t, float64(http.StatusOK), entry["status"])
}

// TestE2E_AuthFailureIsLogged proves a rejected request leaves a trace on
// the server, even though the client only ever sees a bare 401 - closing the
// gap where Clerk rejections used to be entirely silent server-side.
func TestE2E_AuthFailureIsLogged(t *testing.T) {
	buf := captureSlog(t)

	server := newE2ETestServer(t)

	res, err := server.Client().Get(server.URL + "/notes")
	require.NoError(t, err)
	defer res.Body.Close()
	require.Equal(t, http.StatusUnauthorized, res.StatusCode)

	require.Contains(t, buf.String(), "authentication failed")
}
