package integration

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// readBody reads and returns the full response body as a string. The
// response is not closed by readBody - callers keep their existing
// defer/Cleanup for that.
func readBody(t *testing.T, res *http.Response) string {
	t.Helper()

	data, err := io.ReadAll(res.Body)
	require.NoError(t, err)

	return string(data)
}
