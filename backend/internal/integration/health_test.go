package integration

import (
	"net/http"
	"testing"

	"github.com/pitercoding/mindk-ai/backend/internal/testutil"
	"github.com/stretchr/testify/require"
)

// TestHealth_OKWithSecurityAndTracingHeaders proves /health is reachable
// without authentication and that the full middleware stack (security
// headers, request ID) runs in front of it exactly as it does for every
// other route.
func TestHealth_OKWithSecurityAndTracingHeaders(t *testing.T) {
	server := testutil.NewServer(t)

	res, err := server.Client().Get(server.URL + "/health")
	require.NoError(t, err)
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode)

	require.NotEmpty(t, res.Header.Get("X-Request-ID"))
	require.Equal(t, "nosniff", res.Header.Get("X-Content-Type-Options"))
	require.Equal(t, "DENY", res.Header.Get("X-Frame-Options"))
	require.Equal(t, "no-referrer", res.Header.Get("Referrer-Policy"))
	require.Equal(t, "default-src 'none'", res.Header.Get("Content-Security-Policy"))
	require.Equal(t, "no-store", res.Header.Get("Cache-Control"))
}
