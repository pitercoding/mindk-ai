package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHealthHandler(t *testing.T) {

	// Arrange
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	// Act
	HealthHandler(rec, req)

	// Assert
	require.Equal(t, http.StatusOK, rec.Code)

	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var response map[string]string

	err := json.Unmarshal(rec.Body.Bytes(), &response)

	require.NoError(t, err)

	assert.Equal(t, "ok", response["status"])
}