package routes_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pitercoding/mindk-ai/backend/internal/app"
	"github.com/pitercoding/mindk-ai/backend/internal/config"
	"github.com/pitercoding/mindk-ai/backend/internal/routes"
)

// passthrough is a no-op auth middleware, standing in for the real Clerk
// middleware so RegisterRoutes can wire every route without needing a live
// Clerk backend. These tests only ever hit /swagger/, which passthrough
// never even wraps.
func passthrough(next http.Handler) http.Handler {
	return next
}

// TestRegisterRoutes_SwaggerOnlyInDevelopment proves the Swagger UI is wired
// up in development and completely absent (a plain 404, not a 401) in
// production, matching the roadmap 7.9 decision to never expose API docs
// publicly on Render.
func TestRegisterRoutes_SwaggerOnlyInDevelopment(t *testing.T) {
	t.Run("development", func(t *testing.T) {
		mux := http.NewServeMux()
		routes.RegisterRoutes(mux, &app.App{
			Environment:    config.EnvDevelopment,
			AuthMiddleware: passthrough,
		})

		req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusOK, recorder.Code)
	})

	t.Run("production", func(t *testing.T) {
		mux := http.NewServeMux()
		routes.RegisterRoutes(mux, &app.App{
			Environment:    config.EnvProduction,
			AuthMiddleware: passthrough,
		})

		req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
		recorder := httptest.NewRecorder()
		mux.ServeHTTP(recorder, req)

		require.Equal(t, http.StatusNotFound, recorder.Code)
	})
}
