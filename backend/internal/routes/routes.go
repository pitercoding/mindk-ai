package routes

import (
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/app"
	"github.com/pitercoding/mindk-ai/backend/internal/handlers"
)

func RegisterRoutes(app *app.App) {
	// Routes
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/notes", app.NoteHandler.CreateNote)
}
