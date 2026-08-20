package middleware

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/pitercoding/mindk-ai/backend/internal/auth"
	"github.com/stretchr/testify/require"
)

// startMockJWKS serves a JWKS response exposing the public half of the
// given RSA key under the given key ID, mimicking Clerk's JWKS endpoint.
func startMockJWKS(t *testing.T, kid string, pub *rsa.PublicKey) *httptest.Server {
	t.Helper()

	n := base64.RawURLEncoding.EncodeToString(pub.N.Bytes())
	e := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/jwks" && r.Method == http.MethodGet {
			_, err := fmt.Fprintf(w, `{"keys":[{"use":"sig","kty":"RSA","kid":"%s","alg":"RS256","n":"%s","e":"%s"}]}`, kid, n, e)
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

	return server
}

// newProtectedServer wraps a handler that echoes the authenticated user
// ID with the Clerk auth middleware, so tests can hit it over HTTP.
func newProtectedServer(t *testing.T) *httptest.Server {
	t.Helper()

	auth := NewClerkAuth("test_secret_key")

	handler := auth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, ok := clerk.SessionClaimsFromContext(r.Context())
		require.True(t, ok)

		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(map[string]string{"user_id": claims.Subject})
		require.NoError(t, err)
	}))

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return server
}

func TestClerkAuth_NoToken(t *testing.T) {
	server := newProtectedServer(t)

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestClerkAuth_InvalidToken(t *testing.T) {
	server := newProtectedServer(t)

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer invalid-token")

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

func TestClerkAuth_ValidToken(t *testing.T) {
	kid := "kid-" + t.Name()

	token, pub := clerktest.GenerateJWT(t, map[string]any{
		"sub": "user_2abc123",
		"sid": "sess_123",
		"iss": "https://clerk.test.dev",
	}, kid)

	rsaPub, ok := pub.(*rsa.PublicKey)
	require.True(t, ok)

	startMockJWKS(t, kid, rsaPub)

	server := newProtectedServer(t)

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	var body map[string]string
	require.NoError(t, json.NewDecoder(res.Body).Decode(&body))
	require.Equal(t, "user_2abc123", body["user_id"])
}

// TestClerkAuth_EmptySubject proves that a structurally valid token which
// carries no usable identity (empty "sub") is rejected as unauthenticated,
// rather than letting an empty userID flow into the context and downstream
// into repository ownership checks.
func TestClerkAuth_EmptySubject(t *testing.T) {
	kid := "kid-" + t.Name()

	token, pub := clerktest.GenerateJWT(t, map[string]any{
		"sub": "",
		"sid": "sess_123",
		"iss": "https://clerk.test.dev",
	}, kid)

	rsaPub, ok := pub.(*rsa.PublicKey)
	require.True(t, ok)

	startMockJWKS(t, kid, rsaPub)

	server := newProtectedServer(t)

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusUnauthorized, res.StatusCode)
}

// TestClerkAuth_ContextBridge proves the full 6C/6D bridge: a valid Clerk
// JWT results in claims.Subject being placed into the request context via
// auth.WithUserID, and retrievable downstream via auth.UserIDFromContext
// exactly as handlers do - not just via the raw Clerk claims.
func TestClerkAuth_ContextBridge(t *testing.T) {
	kid := "kid-" + t.Name()

	token, pub := clerktest.GenerateJWT(t, map[string]any{
		"sub": "user_2abc123",
		"sid": "sess_123",
		"iss": "https://clerk.test.dev",
	}, kid)

	rsaPub, ok := pub.(*rsa.PublicKey)
	require.True(t, ok)

	startMockJWKS(t, kid, rsaPub)

	authMiddleware := NewClerkAuth("test_secret_key")

	handler := authMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := auth.UserIDFromContext(r.Context())
		require.True(t, ok, "auth.UserIDFromContext must find the userID set by the middleware")

		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(map[string]string{"user_id": userID})
		require.NoError(t, err)
	}))

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	var body map[string]string
	require.NoError(t, json.NewDecoder(res.Body).Decode(&body))
	require.Equal(t, "user_2abc123", body["user_id"])
}
