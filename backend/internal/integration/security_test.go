package integration

import (
	"net/http"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/testutil"
	"github.com/stretchr/testify/require"
)

// TestCORS_AllowedOriginGetsAccessControlHeader proves a request from the
// configured frontend origin receives the matching
// Access-Control-Allow-Origin header, driven through the real CORS
// middleware wired in front of the app (not a unit test of the cors library
// in isolation).
func TestCORS_AllowedOriginGetsAccessControlHeader(t *testing.T) {
	server := testutil.NewServer(t)

	req, err := http.NewRequest(http.MethodGet, server.URL+"/health", nil)
	require.NoError(t, err)
	req.Header.Set("Origin", testutil.TestFrontendOrigin)

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)
	require.Equal(t, testutil.TestFrontendOrigin, res.Header.Get("Access-Control-Allow-Origin"))
}

// TestCORS_DisallowedOriginGetsNoAccessControlHeader proves a request from
// an origin outside FRONTEND_ORIGIN is not granted CORS access: the browser
// enforces same-origin by the absence of the header, since a bare GET can't
// be blocked server-side.
func TestCORS_DisallowedOriginGetsNoAccessControlHeader(t *testing.T) {
	server := testutil.NewServer(t)

	req, err := http.NewRequest(http.MethodGet, server.URL+"/health", nil)
	require.NoError(t, err)
	req.Header.Set("Origin", "https://evil.example")

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Empty(t, res.Header.Get("Access-Control-Allow-Origin"))
}

// TestCORS_PreflightForDisallowedOriginIsRejected proves a real browser
// preflight (OPTIONS with Access-Control-Request-Method) from an
// unauthorized origin is not granted access, which is what actually stops a
// cross-origin browser request from ever reaching a protected endpoint like
// /notes.
func TestCORS_PreflightForDisallowedOriginIsRejected(t *testing.T) {
	server := testutil.NewServer(t)

	req, err := http.NewRequest(http.MethodOptions, server.URL+"/notes", nil)
	require.NoError(t, err)
	req.Header.Set("Origin", "https://evil.example")
	req.Header.Set("Access-Control-Request-Method", http.MethodPost)

	res, err := server.Client().Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	require.Empty(t, res.Header.Get("Access-Control-Allow-Origin"))
}

// TestUnauthenticated_ProtectedEndpointsReject401 sweeps every protected
// route family and proves none of them can be reached without a token,
// closing the gap where a newly added route might forget to wrap itself in
// the auth middleware.
func TestUnauthenticated_ProtectedEndpointsReject401(t *testing.T) {
	server := testutil.NewServer(t)

	endpoints := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/notes"},
		{http.MethodGet, "/documents"},
		{http.MethodGet, "/chat/sessions"},
		{http.MethodGet, "/chat/messages/1"},
		{http.MethodPost, "/chat"},
	}

	for _, ep := range endpoints {
		req, err := http.NewRequest(ep.method, server.URL+ep.path, nil)
		require.NoError(t, err)

		res, err := server.Client().Do(req)
		require.NoError(t, err)
		res.Body.Close()

		require.Equalf(t, http.StatusUnauthorized, res.StatusCode, "%s %s must require authentication", ep.method, ep.path)
	}
}

// TestMalformedInput_RejectedAcrossResources proves invalid JSON is rejected
// with 400 (never reaching the database) for every resource that accepts a
// JSON body, not just notes.
func TestMalformedInput_RejectedAcrossResources(t *testing.T) {
	server := testutil.NewServer(t)
	token := testutil.TokenFor(t, "user_a")

	endpoints := []string{"/notes", "/documents", "/chat/sessions", "/chat"}

	for _, path := range endpoints {
		res := authedRequest(t, server.Client(), token, http.MethodPost, server.URL+path, []byte(`{"title":`))
		require.Equalf(t, http.StatusBadRequest, res.StatusCode, "POST %s with malformed JSON must be rejected with 400", path)
	}
}

// TestInvalidResourceID_Rejected proves a non-numeric path ID is rejected
// with 400 rather than falling through to a 500 from a failed int
// conversion deeper in the stack.
func TestInvalidResourceID_Rejected(t *testing.T) {
	server := testutil.NewServer(t)
	token := testutil.TokenFor(t, "user_a")

	endpoints := []string{"/notes/abc", "/documents/abc", "/chat/sessions/abc"}

	for _, path := range endpoints {
		res := authedRequest(t, server.Client(), token, http.MethodGet, server.URL+path, nil)
		require.Equalf(t, http.StatusBadRequest, res.StatusCode, "GET %s with a non-numeric id must be rejected with 400", path)
	}
}
