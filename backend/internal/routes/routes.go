package routes

import (
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/handlers"
)

func RegisterRoutes() {
	http.HandleFunc("/health", handlers.HealthHandler)
}
