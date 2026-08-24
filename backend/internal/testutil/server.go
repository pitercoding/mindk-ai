// Package testutil builds a fully-wired MindK AI backend - real Clerk auth,
// real security headers, real request logging, real per-user rate limiting
// and real CORS - behind an httptest.Server backed by a temporary, per-test
// SQLite database, so integration tests exercise the exact stack routes.go
// wires in production. The only swapped-out dependency is the LLM: chat
// completions and embeddings are served by a FakeLLMClient so tests never
// call OpenAI.
package testutil

import (
	"crypto/rsa"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/clerktest"
	"github.com/rs/cors"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"

	"github.com/pitercoding/mindk-ai/backend/internal/app"
	"github.com/pitercoding/mindk-ai/backend/internal/config"
	"github.com/pitercoding/mindk-ai/backend/internal/middleware"
	"github.com/pitercoding/mindk-ai/backend/internal/migrations"
	"github.com/pitercoding/mindk-ai/backend/internal/routes"
	"github.com/pitercoding/mindk-ai/backend/internal/services/mocks"
)

// TestFrontendOrigin is the only Origin CORS accepts on a Server built by
// NewServer. It matches the local dev frontend default so CORS tests have a
// realistic allowed value to assert against.
const TestFrontendOrigin = "http://localhost:5173"

// testClerkSecretKey is a placeholder: Clerk verification in tests never
// contacts the real Clerk API, it verifies tokens against a mock JWKS server
// started by TokenFor/TokensFor/StartJWKS.
const testClerkSecretKey = "test_secret_key"

// Server is a full MindK AI backend - every real route, every real
// middleware - running against an isolated, temporary SQLite database.
type Server struct {
	*httptest.Server
	DB  *sql.DB
	LLM *mocks.FakeLLMClient
}

// NewServer builds and starts a Server. The database is a fresh SQLite file
// under t.TempDir(), migrated on startup and discarded when the test ends -
// tests never touch backend/data/mindk.db. The LLM is a FakeLLMClient
// (exposed as Server.LLM) that a test can configure before making requests
// that reach the chat or document-upload paths.
func NewServer(t *testing.T) *Server {
	t.Helper()

	dbPath := filepath.Join(t.TempDir(), "mindk-test.db")

	db, err := sql.Open("sqlite", "file:"+dbPath+"?_busy_timeout=5000")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	db.SetMaxOpenConns(1)
	require.NoError(t, migrations.Run(db))

	cfg := &config.Config{
		Environment:    config.EnvDevelopment,
		OpenAIAPIKey:   "unused-fake-llm-is-injected-instead",
		ClerkSecretKey: testClerkSecretKey,
		FrontendOrigin: TestFrontendOrigin,
		DatabasePath:   dbPath,
	}

	llmClient := &mocks.FakeLLMClient{Response: "test answer"}

	application := app.New(db, cfg, llmClient)

	mux := http.NewServeMux()
	routes.RegisterRoutes(mux, application)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{cfg.FrontendOrigin},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	}).Handler(mux)

	handler := middleware.RequestLogger(middleware.SecurityHeaders(corsHandler))

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return &Server{Server: server, DB: db, LLM: llmClient}
}

// GenerateToken creates a valid but not-yet-verifiable Clerk JWT for userID.
// The caller must start a mock JWKS endpoint for it afterwards - see TokenFor
// for a single live token, or StartJWKS/TokensFor when a test needs several
// users' tokens verifiable at once.
func GenerateToken(t *testing.T, userID string) (token, kid string, pub *rsa.PublicKey) {
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

// TokenFor generates a valid Clerk JWT for userID and starts the mock JWKS
// endpoint needed to verify it. Use this when a test only needs one
// authenticated identity live at a time.
func TokenFor(t *testing.T, userID string) string {
	t.Helper()

	token, kid, pub := GenerateToken(t, userID)
	StartJWKS(t, map[string]*rsa.PublicKey{kid: pub})

	return token
}

// TokensFor generates a valid Clerk JWT for each of userIDs and starts one
// mock JWKS endpoint serving all of their keys, so a test can hold several
// live tokens at once (e.g. proving user A cannot reach user B's data). The
// Clerk backend is a single process-global pointer, so don't call TokenFor
// or StartJWKS again afterwards in the same test - it would replace the JWKS
// endpoint these tokens depend on.
func TokensFor(t *testing.T, userIDs ...string) map[string]string {
	t.Helper()

	tokens := make(map[string]string, len(userIDs))
	keys := make(map[string]*rsa.PublicKey, len(userIDs))

	for _, userID := range userIDs {
		token, kid, pub := GenerateToken(t, userID)
		tokens[userID] = token
		keys[kid] = pub
	}

	StartJWKS(t, keys)

	return tokens
}

// StartJWKS serves the given public keys from one JWKS endpoint and points
// the process-global Clerk backend at it, so Clerk verification succeeds for
// any token signed with the corresponding private key.
func StartJWKS(t *testing.T, keys map[string]*rsa.PublicKey) {
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
