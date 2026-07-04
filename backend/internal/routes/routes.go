package routes

import (
	"net/http"

	"github.com/pitercoding/mindk-ai/backend/internal/database"
	"github.com/pitercoding/mindk-ai/backend/internal/handlers"
	"github.com/pitercoding/mindk-ai/backend/internal/repository"
)

func RegisterRoutes() {
	// Repository
	noteRepo := repository.NewNoteRepository(database.DB)

	// Handler
	noteHandler := handlers.NewNoteHandler(noteRepo)

	// Routes
	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/notes", noteHandler.CreateNote)
}
